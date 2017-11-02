package handler

import (
	"api/model"
	apiHttp "api/http"

	"github.com/labstack/echo"
	"net/http"
)

func CityGetShowcase(c echo.Context) error {

	payload := &apiHttp.Base64Request{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("bind", err.Error())})
	}

	request, err := apiHttp.NewSearchRequest(payload.Search)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("unmarshal", err.Error())})
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("validate", err.Error())})
	}

	cityId := c.Param("cityId")

	return c.JSON(http.StatusOK, &apiHttp.Response{Meta: model.CityNewInventory(cityId, request), Data: model.CityNewComparisonSetsMap(cityId, request)})
}

func CityGetActivities(c echo.Context) error {

	payload := &apiHttp.Base64Request{}
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("bind", err.Error())})
	}

	request, err := apiHttp.NewSearchRequest(payload.Search)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("unmarshal", err.Error())})
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &apiHttp.Response{Data: apiHttp.NewErrorData("validate", err.Error())})
	}

	cityId := c.Param("cityId")

	return c.JSON(http.StatusOK, &apiHttp.Response{Meta: model.CityUpdateInventory(cityId, request), Data: model.CityNewComparisonSetsMap(cityId, request)})
}
