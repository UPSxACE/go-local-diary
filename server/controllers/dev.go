package controllers

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/UPSxACE/go-local-diary/server/dev_component_parser"
	"github.com/labstack/echo/v4"
)

func SetDevRoutes(e *echo.Echo, appConfig *app_config.AppConfig) {
	e.GET("/dev", dev_component_parser.SetDevControllerWrapper(GetDevController, appConfig))
	e.GET("/dev/components", dev_component_parser.SetDevControllerWrapper(GetDevComponentsController, appConfig))
}

func GetDevController(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/dev/components")
}

func GetDevComponentsController(c echo.Context) error {
	renderFunc := dev_component_parser.GetDevComponentParserRenderFunc(c)
	return renderFunc(http.StatusOK, "dev-components")
}


