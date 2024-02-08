package controllers

import (
	"net/http"
	"strconv"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/models"
	"github.com/UPSxACE/go-local-diary/server/internal/services"
	"github.com/UPSxACE/go-local-diary/server/modules/echo_custom"
	"github.com/labstack/echo/v4"
)

type NoteController struct {
	echo *echo.Echo
	app  *app.App[db_sqlite3.Database_Sqlite3]
}

func SetNoteController(e *echo.Echo, appInstance *app.App[db_sqlite3.Database_Sqlite3]) {
	ctrl := &NoteController{echo: e, app: appInstance}
	ctrl.SetRoutes()
}

func (ctrl *NoteController) SetRoutes() {
	welcomeMiddleware := echo_custom.RedirectNotConfiguredToWelcomeMiddleware

	ctrl.echo.GET("/note/:id", welcomeMiddleware(ctrl.getIdRoute()))
	ctrl.echo.GET("/note/:id/edit", welcomeMiddleware(ctrl.getIdEditRoute()))
	ctrl.echo.POST("/note/:id/edit", welcomeMiddleware(ctrl.postIdEditRoute()))
}

func (ctrl *NoteController) getIdRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if(err != nil){
			return err;
		}

		name, err := services.AppConfig.GetName(ctrl.app)
		if err != nil {
			return err
		}


		note, err := services.Note.GetNote(ctrl.app, id)
		if err != nil {
			return err
		}
		contentHTML := models.ParseNoteContentToHTML(note.Content)

		notes, err := services.Note.GetNotesOrderByCreateDateDesc(ctrl.app,"", false)
		if err != err {
			return err
		}

		data := map[string]any{
			"Name": name,
			"Note": note,
			"ContentHTML": contentHTML,
			"Notes": notes,
		}
		return c.Render(http.StatusOK, "note-view", data)
	}
}

func (ctrl *NoteController) getIdEditRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if(err != nil){
			return err;
		}

		name, err := services.AppConfig.GetName(ctrl.app)
		if err != nil {
			return err
		}


		note, err := services.Note.GetNote(ctrl.app, id)
		if err != nil {
			return err
		}

		data := map[string]any{
			"Name": name,
			"Note": note,
		}
		return c.Render(http.StatusOK, "note-edit", data)
	}
}

func (ctrl *NoteController) postIdEditRoute() func(c echo.Context) error {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if(err != nil){
			return err;
		}

		err = c.Request().ParseMultipartForm(1073741824) // 1gb
		if err != nil {
			return err
		}

		ctx := c.Request().Context();

		title := c.FormValue("title")
		content := c.FormValue("content")

		valid, errMsg, err := services.Note.UpdateNote(ctrl.app, ctx,id,title,content)
		if err != nil {
			println(valid, errMsg)
			return err
		}
		// FIXME return validation errors

		return c.Redirect(http.StatusFound, "/note/" + idStr)
	}
}
