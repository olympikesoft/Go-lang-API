package main

import (
	"api/database"
	"api/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"config"
)

func main() {

	database.Connect()

	e := echo.New()

	//client, _ := elastic.NewClient(
	//	elastic.SetURL("http://elastic:changeme@localhost:9200"),
	//)
	//ctx := context.Background()

	e.Use(middleware.Gzip())

	e.GET("/cities/:cityId/showcase", handler.CityGetShowcase)
	e.GET("/cities/:cityId/activities", handler.CityGetActivities)
	//e.Logger.Fatal(e.Start(":1324")) DEVELOPMENT
	e.Logger.Fatal(e.Start(config.GetConfig().Port))
}
//
//
//func respond(c echo.Context, status int, response *SearchResponse) error {
//	return c.JSON(status, response)
//}
