package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetIndexRoutes(e *echo.Echo) {
	e.GET("/", GetDevComponentsController)
}

func GetIndexController(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}