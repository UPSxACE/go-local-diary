package controllers

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/services"
	"github.com/UPSxACE/go-local-diary/server/modules/echo_custom"
	"github.com/labstack/echo/v4"
)

type IndexController struct {
	echo *echo.Echo
	app  *app.App[db_sqlite3.Database_Sqlite3]
}

func SetIndexController(e *echo.Echo, appInstance *app.App[db_sqlite3.Database_Sqlite3]) {
	ctrl := &IndexController{echo: e, app: appInstance}
	ctrl.SetRoutes()
}

func (ctrl *IndexController) SetRoutes() {
	welcomeMiddleware := echo_custom.RedirectNotConfiguredToWelcomeMiddleware

	ctrl.echo.GET("/", welcomeMiddleware(ctrl.getIndexRoute()))
	ctrl.echo.GET("/welcome", welcomeMiddleware(ctrl.getWelcomeRoute()))
	ctrl.echo.POST("/welcome", welcomeMiddleware(ctrl.postWelcomeRoute()))
	ctrl.echo.GET("/404", ctrl.get404Route())
}

func (ctrl *IndexController) getIndexRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	}
}

func (ctrl *IndexController) getWelcomeRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		cc := c.(*echo_custom.CustomEchoContext)

		if cc.IsConfigured {
			return cc.Redirect(http.StatusMovedPermanently, "/")
		}

		return cc.Render(http.StatusOK, "welcome", nil)
	}
}

func (ctrl *IndexController) postWelcomeRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		cc := c.(*echo_custom.CustomEchoContext)

		if cc.IsConfigured {
			return cc.Redirect(http.StatusMovedPermanently, "/")
		}

		ctx := cc.Request().Context();
		
		_, err := services.AppConfig.SetConfiguration(ctrl.app, ctx, "configured", "1")
		if(err != nil){
			return nil
		}

		return cc.Redirect(http.StatusMovedPermanently, "/")
	}

}

func (ctrl *IndexController) get404Route() func(c echo.Context) error {
	return func(c echo.Context) error {

		return c.Render(http.StatusOK, "404", nil)
	}
}
