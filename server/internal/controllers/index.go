package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
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
	ctrl.echo.GET("/new", welcomeMiddleware(ctrl.getNewRoute()))
	ctrl.echo.POST("/new", welcomeMiddleware(ctrl.postNewRoute()))
	ctrl.echo.GET("/404", ctrl.get404Route())
}

func (ctrl *IndexController) getIndexRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		name, err := services.AppConfig.GetName(ctrl.app)
		if err != nil {
			return err
		}

		data := map[string]string{
			"Name": name,
		}

		return c.Render(http.StatusOK, "index", data)
	}
}

func (ctrl *IndexController) getNewRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		name, err := services.AppConfig.GetName(ctrl.app)
		if err != nil {
			return err
		}

		data := map[string]string{
			"Name": name,
		}

		return c.Render(http.StatusOK, "new", data)
	}
}

func (ctrl *IndexController) postNewRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		name, err := services.AppConfig.GetName(ctrl.app)
		if err != nil {
			return err
		}

		data := map[string]string{
			"Name": name,
		}

		ctx := c.Request().Context()

		err = c.Request().ParseMultipartForm(1073741824) // 1gb
		if(err != nil){
			return nil
		}

		newNote := models.NoteModel{
			Title:   c.FormValue("title"),
			Content: c.FormValue("content"),
		}

		valid, newNote, errMsg, err := services.Note.CreateNote(ctrl.app, ctx, newNote)
		if err != nil {
			println(valid, errMsg)
			return err
		}
		// FIXME send back validation errors

		hxObject := map[string]map[string]models.NoteModel{
			"hx:new-note": {
				"data": newNote,
			},
		}

		hxObjectJson, err := json.Marshal(hxObject)
		if err != nil {
			return err
		}

		c.Response().Header().Set("HX-Trigger", string(hxObjectJson))

		return c.Render(http.StatusOK, "new", data)
	}
}

func (ctrl *IndexController) getWelcomeRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		cc := c.(*echo_custom.CustomEchoContext)

		if cc.IsConfigured {
			return cc.Redirect(http.StatusMovedPermanently, "/")
		}

		step := c.QueryParam("step")
		if step == "2" {
			return cc.Render(http.StatusOK, "welcome-step-2", nil)
		}

		return cc.Render(http.StatusOK, "welcome", nil)
	}
}

func (ctrl *IndexController) postWelcomeRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		cc := c.(*echo_custom.CustomEchoContext)

		step := cc.QueryParam("step")
		if cc.IsConfigured || step != "2" {
			return cc.Redirect(http.StatusMovedPermanently, "/")
		}

		ctx := cc.Request().Context()
		name := cc.FormValue("name")
		valid, errMsg, err := services.AppConfig.SetNameConfiguration(ctrl.app, ctx, name)
		if err != nil {
			return err
		}
		if !valid {
			data := map[string]interface{}{
				"form_err": map[string]string{
					"name": errMsg,
				},
			}
			return cc.Render(http.StatusOK, "welcome-step-2", data)
		}

		_, err = services.AppConfig.SetConfiguration(ctrl.app, ctx, "configured", "1")
		if err != nil {
			return nil
		}

		cc.Response().Header().Set("HX-Redirect", "/")
		return cc.NoContent(http.StatusMovedPermanently)
	}

}

func (ctrl *IndexController) get404Route() func(c echo.Context) error {
	return func(c echo.Context) error {

		return c.Render(http.StatusOK, "404", nil)
	}
}
