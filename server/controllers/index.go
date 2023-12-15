package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetIndexRoutes(e *echo.Echo) {
	e.GET("/", GetIndexController)
	e.GET("/404", Get404Controller)
}

func GetIndexController(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func Get404Controller(c echo.Context) error {
	return c.Render(http.StatusOK, "404", nil)
}