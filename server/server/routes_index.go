package server

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/server/models"
	"github.com/UPSxACE/go-local-diary/server/modules/echo_custom"
	"github.com/labstack/echo/v4"
)

func (server *Server) setIndexRoutes() {
	welcomeMiddleware := echo_custom.RedirectNotConfiguredToWelcomeMiddleware

	g := server.Echo.Group("", welcomeMiddleware)
	g.Add("GET", "/", server.getIndexRoute)
	g.Add("GET", "/welcome", server.getWelcomeRoute)
	g.Add("POST", "/welcome", server.postWelcomeRoute)
	g.Add("GET", "/new", server.getNewRoute)
	g.Add("POST", "/new", server.postNewRoute)

	server.Echo.GET("/404", server.get404Route)
}

func (server *Server) getIndexRoute(c echo.Context) error {

	name, err := server.Services.GetUserName()
	if err != nil {
		return err
	}

	search := c.QueryParam("search")

	notes, err := server.Services.GetNotesOrderByCreateDateDesc(search, false)
	if err != err {
		return err
	}

	notesPreview := make([]models.NoteModelPreview, 0, len(notes))
	for _, note := range notes {
		notesPreview = append(notesPreview, models.NewNotePreviewModel(note, 0, 500))
	}

	data := map[string]any{
		"Name":  name,
		"Notes": notesPreview,
	}

	return c.Render(http.StatusOK, "index", data)

}

func (server *Server) getNewRoute(c echo.Context) error {

	name, err := server.Services.GetUserName()
	if err != nil {
		return err
	}

	data := map[string]string{
		"Name": name,
	}

	return c.Render(http.StatusOK, "new", data)

}

func (server *Server) postNewRoute(c echo.Context) error {

	ctx := c.Request().Context()

	err := c.Request().ParseMultipartForm(1073741824) // 1gb
	if err != nil {
		return err
	}

	title := c.FormValue("title")
	content := c.FormValue("content")

	valid, errMsg, err := server.Services.CreateNote(ctx, title, content)
	if err != nil {
		println(valid, errMsg)
		return err
	}
	// FIXME send back validation errors

	return c.Redirect(http.StatusFound, "/")

}

func (server *Server) getWelcomeRoute(c echo.Context) error {

	cc := c.(*echo_custom.CustomEchoContext)

	if cc.IsConfigured {
		isHtmxBoosted := cc.Request().Header.Get("HX-Boosted") != ""

		if isHtmxBoosted {
			cc.Response().Header().Set("HX-Redirect", "/")
			return cc.NoContent(http.StatusOK)
		}

		return cc.Redirect(http.StatusMovedPermanently, "/")
	}

	step := c.QueryParam("step")
	if step == "2" {
		return cc.Render(http.StatusOK, "welcome-step-2", nil)
	}

	return cc.Render(http.StatusOK, "welcome", nil)

}

func (server *Server) postWelcomeRoute(c echo.Context) error {

	cc := c.(*echo_custom.CustomEchoContext)

	step := cc.QueryParam("step")
	if cc.IsConfigured || step != "2" {
		isHtmxBoosted := cc.Request().Header.Get("HX-Boosted") != ""

		if isHtmxBoosted {
			cc.Response().Header().Set("HX-Redirect", "/")
			return cc.NoContent(http.StatusOK)
		}

		return cc.Redirect(http.StatusMovedPermanently, "/")
	}

	ctx := cc.Request().Context()
	name := cc.FormValue("name")
	valid, errMsg, err := server.Services.SetUserName(ctx, name)
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

	_, err = server.Services.SetConfiguration(ctx, "configured", "1")
	if err != nil {
		return err
	}

	cc.Response().Header().Set("HX-Redirect", "/")
	return cc.NoContent(http.StatusMovedPermanently)

}

func (server *Server) get404Route(c echo.Context) error {
	return c.Render(http.StatusOK, "404", nil)
}
