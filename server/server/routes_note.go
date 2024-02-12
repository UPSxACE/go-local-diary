package server

import (
	"net/http"
	"strconv"

	"github.com/UPSxACE/go-local-diary/server/models"
	"github.com/UPSxACE/go-local-diary/server/modules/echo_custom"
	"github.com/labstack/echo/v4"
)

func (server *Server) setNoteRoutes() {
	welcomeMiddleware := echo_custom.RedirectNotConfiguredToWelcomeMiddleware

	g := server.Echo.Group("/note", welcomeMiddleware)
	g.Add("GET", "/:id", server.getIdRoute)
	g.Add("GET", "/:id/edit", server.getIdEditRoute)
	g.Add("POST", "/:id/edit", server.postIdEditRoute)
}

func (server *Server) getIdRoute(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	name, err := server.Services.GetUserName()
	if err != nil {
		return err
	}

	note, err := server.Services.GetNote(id)
	if err != nil {
		return err
	}
	contentHTML := models.ParseNoteContentToHTML(note.Content)

	notes, err := server.Services.GetNotesOrderByCreateDateDesc("", false)
	if err != err {
		return err
	}

	notesPreview := make([]models.NoteModelPreview, 0, len(notes))
	for _, note := range notes {
		notesPreview = append(notesPreview, models.NewNotePreviewModel(note, 0, 250))
	}

	data := map[string]any{
		"Name":        name,
		"Note":        note,
		"ContentHTML": contentHTML,
		"Notes":       notesPreview,
	}
	return c.Render(http.StatusOK, "note-view", data)

}

func (server *Server) getIdEditRoute(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	name, err := server.Services.GetUserName()
	if err != nil {
		return err
	}

	note, err := server.Services.GetNote(id)
	if err != nil {
		return err
	}

	data := map[string]any{
		"Name": name,
		"Note": note,
	}
	return c.Render(http.StatusOK, "note-edit", data)

}

func (server *Server) postIdEditRoute(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	err = c.Request().ParseMultipartForm(1073741824) // 1gb
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	title := c.FormValue("title")
	content := c.FormValue("content")

	valid, errMsg, err := server.Services.UpdateNote(ctx, id, title, content)
	if err != nil {
		println(valid, errMsg)
		return err
	}
	// FIXME return validation errors

	return c.Redirect(http.StatusFound, "/note/"+idStr)

}
