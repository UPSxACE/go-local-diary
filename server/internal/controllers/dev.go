package controllers

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/labstack/echo/v4"
)

type DevController struct {
	echo *echo.Echo
	app  *app.App[db_sqlite3.Database_Sqlite3]
}

func SetDevController(e *echo.Echo, appInstance *app.App[db_sqlite3.Database_Sqlite3]) {
	ctrl := &DevController{echo: e, app: appInstance}
	ctrl.SetRoutes()
}

func (ctrl *DevController) SetRoutes() {
	ctrl.echo.GET("/dev", ctrl.getDevController())
	ctrl.echo.GET("/dev/components", dev_component_parser.SetDevControllerWrapper(ctrl.getDevComponentsController(), ctrl.app))
	ctrl.echo.GET("/dev/components/refresh", dev_component_parser.SetDevComponentsRefreshRoute(ctrl.app))
}

func (ctrl *DevController) getDevController() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/dev/components")
	}
}

func (ctrl *DevController) getDevComponentsController() func(c echo.Context) error {
	return func(c echo.Context) error {
		renderFunc := dev_component_parser.GetDevComponentParserRenderFunc(c)
		return renderFunc(http.StatusOK, "dev-components")
	}
}
