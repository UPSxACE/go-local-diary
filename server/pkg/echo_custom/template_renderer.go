/*
The custom package contains the necessary struct and
Render function required by the echo server and a custom context struct.
*/

package echo_custom

import (
	"fmt"
	"html/template"
	"io"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}
type TemplateDevMode struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {	
	// Ensure data is a map[string]interface{}
	newData, ok := data.(map[string]interface{})
	if !ok {
		// If data is not a map, do nothing
	} else {
		httpOrHttps := c.Scheme()
		newData["HOST"] = fmt.Sprintf("%v://%v", httpOrHttps, c.Request().Host)
	}
	
	
	return t.Templates.ExecuteTemplate(w, name, data)
}

func (t *TemplateDevMode) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// In developer mode, the templates are parsed on each request
	tBuilder := template.Must(template.New("").Funcs(app.DefaultFuncMap).ParseGlob("server/internal/views/*/*.html"))

	// tBuilder = template.Must(tBuilder.ParseGlob("server/internal/views/*/*/*.html"))
	tNew := &Template{
		Templates: tBuilder,
	}

	// Ensure data is a map[string]interface{}
	newData, ok := data.(map[string]interface{})
	if !ok {
		// If data is not a map, do nothing
	} else {
		httpOrHttps := c.Scheme()
		newData["HOST"] = fmt.Sprintf("%v://%v", httpOrHttps, c.Request().Host)
	}
	
	return tNew.Templates.ExecuteTemplate(w, name, data)
}
