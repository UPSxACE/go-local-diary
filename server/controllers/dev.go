package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetDevRoutes(e *echo.Echo) {
	e.GET("/dev", GetDevController)
	e.GET("/dev/components", GetDevComponentsController)
}

func GetDevController(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/dev/components")
}

func GetDevComponentsController(c echo.Context) error {
	return c.Render(http.StatusOK, "dev-components", nil)
}