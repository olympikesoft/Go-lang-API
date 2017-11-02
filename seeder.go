package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"crypto/md5"
	"encoding/hex"

	"regexp"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type ActivityLine struct {
	id            int
	comparisonKey string
	ranking       int
	city          int
	identities    string
	categories    string
	primaryTags   string
	provider      string
	operator	  string
	surrogateId   string
	url           string
	secondaryTags string
	availability  string
	languages     string
}

func main() {

	file, err := os.Open("data/Namings - Comparation tvt.csv")
	if err != nil {
		logrus.Error(err.Error())
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	db, err = sqlx.Open("mysql", "wajanga-dev:wajanga-dev@/wajanga-dev")
	if err != nil {
		logrus.Error(err.Error())
	}
	defer db.Close()

	tvts := []*ActivityLine{}
	hasher := md5.New()

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if line[0] == "Comparation set" {
			continue
		}

		// activities without id yet, ceetiz and expedia
		if line[7] == "" {
			continue
		}

		var slugFormatter = strings.NewReplacer(";", ",", "", "", "'", "", ".", "", " ", "", "\r", "", "\t", " ")
		operator := strings.ToLower(strings.trimSpace(line[6]))
		provider := strings.ToLower(strings.TrimSpace(line[6]))
		ranking, _ := strconv.Atoi(line[1])
		categories := strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[4]), ",")))
		city, _ := strconv.Atoi(line[2])
		activityId := strings.TrimSpace(line[7])

		if activityId == "" {
			exp := regexp.MustCompile(`\?tour=(\d+)\&`)
			match := exp.FindStringSubmatch(line[8])

			if match == nil {
				logrus.Fatalf("missing activityId in %q feed, url is: %+v \n", provider, line[8])
			}
			activityId = match[1]
		}

		hasher.Write([]byte(provider))
		hasher.Write([]byte(provider))
		hasher.Write([]byte(activityId))

		tvt := &ActivityLine{
			ranking:       ranking,
			city:          city,
			identities:    strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[3]), ","))),
			categories:    categories,
			primaryTags:   strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[5]), ","))),
			provider:      provider,
			operator:	   operator,
			surrogateId:   hex.EncodeToString(hasher.Sum(nil)),
			url:           strings.TrimSpace(formatAffiliate(provider, line[8])),
			secondaryTags: strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[9]), ","))),
			availability:  strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[11]), ","))),
			languages:     strings.ToLower(strings.TrimSpace(strings.TrimRight(slugFormatter.Replace(line[12]), ","))),
		}

		hasher.Reset()

		hasher.Write([]byte(tvt.identities))
		hasher.Write([]byte(tvt.primaryTags))
		tvt.comparisonKey = hex.EncodeToString(hasher.Sum(nil))
		hasher.Reset()

		tvt.save(line, tvts);
	}

	fmt.Printf("-- %+v \n", len(tvts))
}

func formatAffiliate(provider, url string) string {
	switch provider {
	case "tiqets":
		url = url + "&affiliateID=4880LOCAL"
		break
	case "citydiscovery":
		url = url + "&partner=wajanga"
		url = url + "&partner=wajanga"
		break
	case "musement":
		url = url + "?aid=Wajanga"
		break
	case "getyourguide":
		url = url + "?partner_id=RM45QCZ&psrc=partner_api&currency=EUR"
		break
	case "headout":
		url = "https://headout.go2cloud.org/aff_c?offer_id=4&aff_id=1030&file_id=2&url=" + url
		break
	case "ceetiz":
		url = url
		break
	case "viator":
		url = url //"http://www.partner.viator.com/en/23271"
		break
	case "expedia":
		url = url
		break
	}

	return url
}

func (tvt *ActivityLine) save(line []string, tvts []*ActivityLine) {
	sql_procedure := "CALL insert_activity(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,@idts)"

	stmt, err := db.Prepare(sql_procedure)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Query(
		tvt.city,
		tv.operator,
		tvt.provider,
		tvt.surrogateId,
		tvt.comparisonKey,
		"",
		tvt.url,
		nil,
		"",
		tvt.ranking,
		tvt.categories,
		tvt.identities,
		tvt.primaryTags,
		tvt.secondaryTags,
		tvt.availability,
		tvt.languages,
		nil,
		"",
	)
	if err != nil {
		fmt.Println(tvt.id)
		logrus.Error(err.Error())
	}

	tvts = append(tvts, tvt)
}