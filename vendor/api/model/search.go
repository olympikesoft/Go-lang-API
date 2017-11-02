package model

import elastic "gopkg.in/olivere/elastic.v3"

type SearchResponse struct {
	Request *SearchRequest        `json:"request"`
	Data    *elastic.SearchResult `json:"data"`
}

type SearchRequest struct {
	Match      Match      `json:"match" validate:"required"`
	Categories []int      `json:"activities" form:"activities" query:"activities" validate:"required"`
	Duration   int32      `json:"duration" form:"duration" query:"duration" validate:""`
	PriceRange PriceRange `json:"price_range" validate:"required"`
}

type Match struct {
	Must    string `json:"must" query:"must" validate:"required"`
	MustNot string `json:"mustnot" query:"mustnot" validate:""`
}

type PriceRange struct {
	Min int `json:"min" query:"min" validate:"required"`
	Max int `json:"max" query:"max" validate:"required"`
}
