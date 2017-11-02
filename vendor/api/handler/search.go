package handler

//e.Any("/cities/:city", func(c echo.Context) error {
//	sr := new(SearchRequest)
//	if err := c.Bind(sr); err != nil {
//		logrus.Error(err)
//		return respond(c, http.StatusUnprocessableEntity, &SearchResponse{sr, &elastic.SearchResult{}})
//	}
//
//	if err := c.Validate(sr); err != nil {
//		logrus.Error(err)
//		return respond(c, http.StatusUnprocessableEntity, &SearchResponse{sr, &elastic.SearchResult{}})
//	}
//
//	mainBool := elastic.NewBoolQuery()
//	query := mainBool.Must(elastic.NewMatchQuery("name", sr.Match.Must))
//
//	if sr.Match.MustNot != "" {
//		query = query.MustNot(elastic.NewMatchQuery("name", sr.Match.MustNot))
//	}
//
//	filterBool := elastic.NewBoolQuery()
//	query = query.Filter(filterBool.Must(elastic.NewTermsQuery("activities", sr.Categories[0])))
//
//	if sr.Duration < -10 {
//		query = query.Filter(filterBool.Must(elastic.NewTermQuery("duration", sr.Duration)))
//	}
//
//	rangeQuery := elastic.NewRangeQuery("price.eur")
//	rangeQuery = rangeQuery.Gte(sr.PriceRange.Min).Lte(sr.PriceRange.Max)
//	query = query.Filter(filterBool.Must(rangeQuery))
//
//	topHits := elastic.NewTopHitsAggregation().Sort("_score", false).Size(1)
//
//	comparison := elastic.NewTermsAggregation().Field("provider.keyword")
//	comparison = comparison.SubAggregation("best_match", topHits)
//
//	searchResult, err := client.Search().
//		Index("italy").
//		Type(c.Param("city")).
//		Query(query).
//		Sort("price.eur", true).
//		Size(0).
//		Aggregation("min_price", elastic.NewMinAggregation().Field("price.eur")).
//		Aggregation("max_price", elastic.NewMaxAggregation().Field("price.eur")).
//		Aggregation("comparison", comparison).
//		Do(ctx)
//	if err != nil {
//		// Handle error
//		panic(err)
//	}
//
//	return respond(c, http.StatusOK, &SearchResponse{Request: sr, Data: searchResult})
//})
