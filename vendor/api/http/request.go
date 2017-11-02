package http

import (
	"encoding/json"
	"encoding/base64"
	"github.com/asaskevich/govalidator"
	"errors"
)

type Base64Request struct {
	Search string  `json:"search" query:"search"`
}

type SearchRequest struct {
	Filters     struct {
			   SortPrice string `json:"sort_price"`
		   } `json:"filters"`
	Categories Categories `json:"categories"`
	PriceRange []int `json:"minmax"`
	Page       int `json:"page"`
}

type Categories map[int]struct {
	Idt []string `json:"idt"`
	Tag []string `json:"tag"`
}

func NewSearchRequest(payload string) (*SearchRequest, error) {
	r := &SearchRequest{}

	s, _ := base64.StdEncoding.DecodeString(payload)

	err := json.Unmarshal(s, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (r *SearchRequest) Validate() error {

	if !govalidator.IsPositive(float64(r.Page)) {
		return errors.New("page is not valid")
	}

	if govalidator.Matches(r.Filters.SortPrice, "^[asc|desc]$") {
		return errors.New("sort price filter is not valid")
	}

	return nil

}