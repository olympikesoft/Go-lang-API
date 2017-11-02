package model

import (
	"api/database"
	"api/helpers"
	"api/http"
	"fmt"
	"log"
	"strings"

	"database/sql"

	"github.com/Sirupsen/logrus"
)

const RECORD_PER_PAGE = 10

type ComparisonSet struct {
	Winner      *Activity   `json:"winner"`
	Competitors []*Activity `json:"competitors"`
}

type Activity struct {
	Description  string                       `json:"description"`
	Duration     int64                        `json:"duration"`
	Identities   map[string]string            `json:"identities"`
	Name         string                       `json:"name"`
	Price        uint32                       `json:"price"`
	PrimaryTags  map[string]string            `json:"primary_tags"`
	Provider     string                       `json:"provider"`
	Operator     string                       `json:"operator"`
	ReviewsAvg   float64                      `json:"reviews_avg"`
	ReviewsCount float64                      `json:"reviews_count"`
	Tags         map[string]map[string]string `json:"tags"`
	Thumb        string                       `json:"thumb"`
	Url          string                       `json:"url"`
}

type QueryActivity struct {
	ComparisonKey string `db:"comparison_key"`
	Description   sql.NullString
	Duration      sql.NullInt64
	Id            int
	Identities    string `db:"identities"`
	Name          string
	Price         uint32
	Provider      string
	Operator      string
	Ranking       int
	PrimaryTags   string          `db:"primary_tags"`
	ReviewsAvg    sql.NullFloat64 `db:"reviews_avg"`
	ReviewsCount  sql.NullFloat64 `db:"reviews_count"`
	TagId         *string         `db:"tag_id"`
	TagName       *string         `db:"tag_name"`
	TagType       *string         `db:"tag_type"`
	Thumb         sql.NullString
	Url           string
}

type QueryInventory struct {
	//Id       int `db:"category_id"`
	//Name     string `db:"category_name"`
	//Priority string `db:"category_priority"`
	Claim string `db:"claim"`
	Count int
	Min   uint32 `db:"min"`
	Max   uint32 `db:"max"`
}

type CategoryFilter struct {
	Id   int            `json:"id"`
	Name string         `json:"name"`
	Idts map[int]string `json:"idts"`
	Tags map[int]string `json:"tags"`
}

type Inventory struct {
	Claim       string           `json:"claim"`
	Categories  []CategoryFilter `json:"categories"`
	PricesRange [2]uint32        `json:"minmax"`
	Totals      Totals           `json:"totals"`
}

type Totals struct {
	ComparisonSet int `json:"csets" db:"csets"`
	Activities    int `json:"tvts" db:"tvts"`
}

func CityNewInventory(cityId string, request *http.SearchRequest) *Inventory {

	categoriesFilters := buildCategoriesFilters(cityId)

	inventory := &Inventory{
		Categories: categoriesFilters,
	}

	addPricesAndClaim(cityId, "EUR", inventory)

	inventory.Totals = getTotals(cityId, "EUR", &request.Categories)

	return inventory
}

func CityUpdateInventory(cityId string, request *http.SearchRequest) *Inventory {

	inventory := &Inventory{}

	addPricesAndClaim(cityId, "EUR", inventory)

	inventory.Totals = getTotals(cityId, "EUR", &request.Categories)

	return inventory
}

func CityNewComparisonSetsMap(cityId string, request *http.SearchRequest) map[string]*ComparisonSet {

	var lastActivity int
	vararg := map[string]interface{}{"rpp": RECORD_PER_PAGE, "offset": (request.Page - 1) * 10, "cityId": cityId, "curr": "EUR"}

	activities := fetchActivities(request, vararg)

	sets := make(map[string]*ComparisonSet)
	tags := make(map[int]map[string]map[string]string)
	activitiesAdded := make(map[int]struct{})

	var comparisonKeys []string

	for _, activity_db := range activities {

		isWinner := false

		currentActivity := activity_db.Id

		if activity_db.TagId != nil && activity_db.TagName != nil {
			if lastActivity == 0 || lastActivity != currentActivity {
				tags[activity_db.Id] = map[string]map[string]string{}
			}

			if _, ok := tags[activity_db.Id][*activity_db.TagType]; !ok {
				tagKeyVal := map[string]string{*activity_db.TagId: *activity_db.TagName}
				tags[activity_db.Id][*activity_db.TagType] = tagKeyVal
			} else {
				tags[activity_db.Id][*activity_db.TagType][*activity_db.TagId] = *activity_db.TagName
			}
		}

		if _, ok := activitiesAdded[currentActivity]; ok {
			continue
		}

		lastActivity = currentActivity

		if indexOf(activity_db.ComparisonKey, comparisonKeys) == -1 {
			comparisonKeys = append(comparisonKeys, activity_db.ComparisonKey)
			isWinner = true
		}

		ranking := indexOf(activity_db.ComparisonKey, comparisonKeys)

		currentComparisonKey := fmt.Sprintf("%02d-%x", ranking, activity_db.ComparisonKey)

		identities := make(map[string]string)
		for _, p := range strings.Split(activity_db.Identities, ",[0-9]") {
			t := strings.Split(p, "-")
			identities[t[0]] = t[1]
		}

		primaryTags := make(map[string]string)
		for _, p := range strings.Split(activity_db.PrimaryTags, ",") {
			t := strings.Split(p, "-")
			primaryTags[t[0]] = t[1]
		}

		activity := &Activity{
			Description:  activity_db.Description.String,
			Duration:     activity_db.Duration.Int64,
			Identities:   identities,
			Name:         activity_db.Name,
			Provider:     activity_db.Provider,
			Operator:     activity_db.Operator,
			Price:        helpers.RoundPrice(activity_db.Price),
			PrimaryTags:  primaryTags,
			ReviewsAvg:   activity_db.ReviewsAvg.Float64,
			ReviewsCount: activity_db.ReviewsCount.Float64,
			Tags:         tags[activity_db.Id],
			Thumb:        activity_db.Thumb.String,
			Url:          activity_db.Url,
		}

		if isWinner {
			sets[currentComparisonKey] = &ComparisonSet{Winner: activity}
		} else {
			sets[currentComparisonKey].Competitors = append(sets[currentComparisonKey].Competitors, activity)
		}

		activitiesAdded[currentActivity] = struct{}{}
	}

	return sets
}

func fetchActivities(request *http.SearchRequest, vararg map[string]interface{}) []QueryActivity {

	var activities []QueryActivity

	sql_activities := `
		SELECT a.id, a.comparison_key, COALESCE(a.name, '') as name, a.description, a.url,
		a.duration, a.reviews_avg, a.reviews_count, a.thumb, pr.name provider, p.price,
		oper.name as operator,
		(SELECT GROUP_CONCAT(DISTINCT i.id,'-',IFNULL(i.name, i.label)) FROM activities_identities ai JOIN identities i ON i.id = ai.identity_id WHERE ai.activity_id = a.id) identities,
		(SELECT GROUP_CONCAT(DISTINCT pt2.id,'-',IFNULL(pt2.name, pt2.label)) FROM activities_primary_tags apt JOIN primary_tags pt2 ON pt2.id = apt.primary_tag_id WHERE apt.activity_id = a.id) primary_tags,
		t.id tag_id, t.name tag_name, tt.type tag_type
		FROM activities a
		JOIN (SELECT comparison_key, MIN(p.price), IF(ac.ranking > 0, ac.ranking, (SELECT ranking FROM activities ORDER BY ranking desc LIMIT 1)+1) ranking
		FROM activities ac
		JOIN activities_categories aca ON aca.activity_id = ac.id
		JOIN categories cat ON aca.category_id = cat.id
		JOIN activities_identities ai ON ac.id = ai.activity_id
		JOIN identities i ON i.id = ai.identity_id
		JOIN cities ci ON ac.city_id = ci.id
		JOIN pricetables p ON ac.id = p.activity_id
		LEFT JOIN activities_primary_tags apt ON ac.id = apt.activity_id
		LEFT JOIN primary_tags pt ON pt.id = apt.primary_tag_id
		WHERE ci.id = :cityId AND p.currency = :curr
		`

	sql_activities += composeFromCategories(&request.Categories)

	if len(request.PriceRange) == 2 {
		sql_activities += " AND p.price BETWEEN :min AND :max "
		vararg["min"] = request.PriceRange[0]
		vararg["max"] = request.PriceRange[1]
	}

	sql_activities += " GROUP BY comparison_key, ac.ranking "

	if request.Filters.SortPrice != "" {
		sql_activities += fmt.Sprintf("ORDER BY MIN(p.price) %s ", request.Filters.SortPrice)
		//sortBy = "price"
	} else {
		sql_activities += "ORDER BY ranking "
	}

	sql_activities += `
		LIMIT :rpp OFFSET :offset
		) as ck ON ck.comparison_key = a.comparison_key
		JOIN providers pr ON a.provider_id = pr.id
		LEFT JOIN operators oper ON a.operator_id = oper.id
		JOIN pricetables p ON a.id = p.activity_id
		JOIN activities_primary_tags apt ON a.id = apt.activity_id
		JOIN primary_tags pt ON pt.id = apt.primary_tag_id
		LEFT JOIN activities_tags at ON a.id = at.activity_id
		LEFT JOIN tags t ON t.id = at.tag_id
		LEFT JOIN tags_types tt ON tt.id = t.type_id
		WHERE p.price IS NOT NULL
		`

	sql_activities += "GROUP BY a.id, p.price, t.id, tt.type, ck.ranking "

	if request.Filters.SortPrice != "" {
		sql_activities += fmt.Sprintf("ORDER BY p.price %s, a.ranking, ", request.Filters.SortPrice)
	} else {
		sql_activities += "ORDER BY ck.ranking, p.price, "
	}

	sql_activities += "pr.ranking, a.duration"

	rows, err := database.Mysql.NamedQuery(sql_activities, vararg)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		activity := &QueryActivity{}
		err := rows.StructScan(&activity)
		if err != nil {
			log.Fatalln(err)
		}
		activities = append(activities, *activity)
	}

	return activities
}

func composeFromCategories(categories *http.Categories) string {

	if len(*categories) == 0 {
		return ""
	}

	sql_part := "AND ("
	catCtr := 0

	for catId, cat := range *categories {
		if catCtr > 0 && catCtr <= len(*categories)-1 {
			sql_part += fmt.Sprintf(" OR (cat.id = %v ", catId)
		} else if catCtr == 0 {
			sql_part += fmt.Sprintf(" (cat.id = %v ", catId)
		} else {
			sql_part += fmt.Sprintf(" AND (cat.id = %v ", catId)
		}

		if len(cat.Idt) > 0 {
			sql_part += " AND ( "
			for catIdx, i := range cat.Idt {
				if catIdx > 0 && catIdx <= len(cat.Idt)-1 {
					sql_part += " OR "
				}
				sql_part += fmt.Sprintf(" i.id = %v", i)
			}
			sql_part += " ) "
		}

		if len(cat.Tag) > 0 {
			sql_part += " AND ( "
			for tagIdx, t := range cat.Tag {
				if tagIdx > 0 && tagIdx <= len(cat.Tag)-1 {
					sql_part += " OR "
				}
				sql_part += fmt.Sprintf(" pt.id = %v", t)
			}
			sql_part += " ) "
		}

		sql_part += ")"

		catCtr++
	}

	sql_part += ") "

	return sql_part
}

func buildCategoriesFilters(cityId string) []CategoryFilter {

	sql_catFilter := `
		SELECT c.ui_priority, c.name category_name, c.id category_id, idts.id id, idts.name name, ui.type
		FROM ui_filters ui
		JOIN cities ci ON ui.city_id = ci.id
		JOIN categories c ON c.id = ui.category_id
		JOIN activities_categories ac ON ac.category_id = c.id
		JOIN activities a ON a.id = ac.activity_id
		JOIN identities idts ON idts.id = ui.id_value AND ui.type = 'idt'
		JOIN pricetables p ON a.id = p.activity_id
		AND ci.id = ?
		GROUP BY c.id, idts.id, ui.type
		UNION
		SELECT c.ui_priority, c.name category_name, c.id category_id, ptags.id id, ptags.name name, ui.type
		FROM ui_filters ui
		JOIN cities ci ON ui.city_id = ci.id
		JOIN categories c ON c.id = ui.category_id
		JOIN activities_categories ac ON ac.category_id = c.id
		JOIN activities a ON a.id = ac.activity_id
		JOIN primary_tags ptags ON ptags.id = ui.id_value AND ui.type = 'tag'
		JOIN pricetables p ON a.id = p.activity_id
		AND ci.id = ?
		GROUP BY c.id, ptags.id, ui.type
		ORDER BY ui_priority
		`

	rows, err := database.Mysql.Query(sql_catFilter, cityId, cityId)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer rows.Close()

	catFilters := []CategoryFilter{}

	for rows.Next() {

		var ui_priority int
		var category_name string
		var category_id int
		var id int
		var name string
		var filter_type string

		err := rows.Scan(&ui_priority, &category_name, &category_id, &id, &name, &filter_type)
		if err != nil {
			log.Fatalln(err)
		}

		catFilter, found := has(catFilters, category_id)

		if !found {
			catFilter := CategoryFilter{
				Id:   category_id,
				Name: category_name,
				Idts: map[int]string{},
				Tags: map[int]string{},
			}

			if filter_type == "idt" {
				catFilter.Idts[id] = name
			} else {
				catFilter.Tags[id] = name
			}

			catFilters = append(catFilters, catFilter)
		} else {
			if filter_type == "idt" {
				catFilter.Idts[id] = name
			} else {
				catFilter.Tags[id] = name
			}
		}
	}

	return catFilters
}

func addPricesAndClaim(cityId string, lang string, inventory *Inventory) {
	var query_pr []QueryInventory

	//sql_pr := "SELECT MIN(p.price) min, MAX(p.price) max, c.name category_name, c.id category_id, c.ui_priority category_priority FROM activities a JOIN cities ci ON a.city_id = ci.id JOIN pricetables p ON a.id = p.activity_id JOIN activities_categories ac ON a.id = ac.activity_id JOIN categories c ON c.id = ac.category_id WHERE ci.id = ? AND p.currency = ? GROUP BY c.id ORDER BY category_priority"
	sql_pr := "SELECT ci.claim, MIN(p.price) min, MAX(p.price) max FROM activities a JOIN cities ci ON a.city_id = ci.id JOIN pricetables p ON a.id = p.activity_id JOIN activities_categories ac ON a.id = ac.activity_id JOIN categories c ON c.id = ac.category_id WHERE ci.id = ? AND p.currency = ? GROUP BY c.id, ci.claim"
	err := database.Mysql.Select(&query_pr, sql_pr, cityId, lang)
	if err != nil {
		logrus.Error(err.Error())
	}

	for _, item := range query_pr {

		inventory.Claim = item.Claim

		//cat := &Category{
		//	Id:   item.CategoryId,
		//	Name: item.CategoryName,
		//}

		//inventory.Categories[item.CategoryPriority] = *cat
		if inventory.PricesRange[0] == 0 {
			inventory.PricesRange[0] = item.Min
		}
		inventory.PricesRange[0] = helpers.Min(inventory.PricesRange[0], item.Min)
		inventory.PricesRange[1] = helpers.Max(inventory.PricesRange[1], item.Max)
	}
}

func getTotals(cityId string, lang string, categories *http.Categories) Totals {
	count_activities := `SELECT
		COUNT(DISTINCT a.comparison_key) csets,
		COUNT(DISTINCT a.id) tvts
		FROM activities a
		JOIN activities_categories aca ON aca.activity_id = a.id
		JOIN categories cat ON aca.category_id = cat.id
		JOIN activities_identities ai ON a.id = ai.activity_id
		JOIN identities i ON i.id = ai.identity_id
		JOIN cities ci ON a.city_id = ci.id
		JOIN providers pro ON a.provider_id = pro.id
		JOIN operators oper ON a.operator_id = oper.id
		JOIN pricetables p ON a.id = p.activity_id
		LEFT JOIN activities_primary_tags apt ON a.id = apt.activity_id
		LEFT JOIN primary_tags pt ON pt.id = apt.primary_tag_id
		WHERE ci.id = ? AND p.currency = ?
		AND p.price IS NOT NULL
		`

	count_activities += composeFromCategories(categories)

	var totals Totals

	err := database.Mysql.Get(&totals, count_activities, cityId, lang)
	if err != nil {
		logrus.Error(err.Error())
	}

	return totals
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func has(list []CategoryFilter, item int) (CategoryFilter, bool) {
	for _, v := range list {
		if v.Id == item {
			return v, true
		}
	}

	return CategoryFilter{}, false
}
