package controllers

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/labstack/echo/v4"
)

func SetDevRoutes(e *echo.Echo, app *app.App[db_sqlite3.Database_Sqlite3]) {
	e.GET("/dev", GetDevController)
	e.GET("/dev/components", dev_component_parser.SetDevControllerWrapper(GetDevComponentsController, app))
	e.GET("/dev/components/refresh", dev_component_parser.SetDevComponentsRefreshRoute(app))
}

func GetDevController(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/dev/components")
}

func GetDevComponentsController(c echo.Context) error {
	renderFunc := dev_component_parser.GetDevComponentParserRenderFunc(c)
	return renderFunc(http.StatusOK, "dev-components")
}


