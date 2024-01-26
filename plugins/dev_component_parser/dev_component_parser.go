// Web

/*
The package dev_component_parser is a plugin that can be used
to test templates(including testing them with different data)

To use it a json file("dev-components-json") must be created
with the right format(so it is parsed by the plugin and converted
into a []Category variable), a controller route must be
created with its render method wrapped by GetDevComponentParserRenderFunc,
then that controller must be wrapped by the method SetDevControllerWrapper
in the SetRoute method of the controller, and then
another controller route must be created using the
returned value of SetDevComponentsRefreshRoute.

The correct templates/views must also be defined:
- dev-components
- dev-components-render
- dev-components-showcase-header
- dev-components-side
- dev-components-showcase-sideinfo
*/
package dev_component_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/server/echo_custom"
	"github.com/UPSxACE/go-local-diary/utils"

	"github.com/labstack/echo/v4"
)

/* File that will be parsed to load the examples. */
var jsonPath string = "./server/dev-components.json";

/* The main struct that will be stored in App.Plugins */
type DevComponentParser struct {
	Data []Category
}

/* Constructor. */
func Init() *DevComponentParser {
	devCompPars := DevComponentParser{}
	devCompPars.ParseJsonConfigFile()
	return &devCompPars
}

/* Function used to initialize the plugin. */
func LoadPlugin[T any](app *app.App[T]) {
	fmt.Println("Loading Plugin DevComponentParser...")
	if app.DevMode {
		app.Plugins["DevComponentParser"] = Init()
	}
}

/*
Parses data from the JSON file. This will be automatically
used as soon as the struct is initialized using the Init()
constructor, and then recalled to refresh the data whenever
the refresh route set by SetDevComponentsRefreshRoute is
accessed.
*/
func (devComponentParser *DevComponentParser) ParseJsonConfigFile() *DevComponentParser {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		devComponentParser.Data = []Category{}
		return devComponentParser
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var parsedData []Category

	err = json.Unmarshal([]byte(byteValue), &parsedData)
	if err != nil {
		devComponentParser.Data = []Category{}
		return devComponentParser
	}

	devComponentParser.Data = parsedData

	return devComponentParser
}

/*
High order function to wrap a controller to inject the data
parsed from the JSON file in the context. 
*/
func SetDevControllerWrapper[T any](controller echo.HandlerFunc, app *app.App[T]) echo.HandlerFunc {
	parser, ok := app.Plugins["DevComponentParser"]
	if ok {
		parserConverted, conversionOk := parser.(*DevComponentParser)
		if conversionOk {
			wrappedFunc := func(c echo.Context) error {
				cc := &devContextWrapper{Context: c, DevCompParser: parserConverted}
				return controller(cc)
			}
			return wrappedFunc
		}
	}
	return controller
}

/*
Struct to represent the expected data that must be in the
request.
*/
type DevComponentParserRequest struct {
	Render    string `query:"render"`
	Category  int    `query:"ct"` // (index)
	Component int    `query:"cp"` // (index+1)
	Example   int    `query:"e"`  // (index+1)
}


/*
Struct to extend the echo.Context fields
*/
type devContextWrapper struct {
	echo.Context
	DevCompParser *DevComponentParser
}

/*
High order function to wrap a controller to inject the data
parsed from the JSON file in the context.
*/
func GetDevComponentParserRenderFunc(c echo.Context) func(code int, name string) error {
	devContext, ok := c.(*devContextWrapper)

	if ok {
		return func(code int, name string) error {
			requestData := new(DevComponentParserRequest)
			if err := c.Bind(requestData); err != nil {
				return c.String(http.StatusBadRequest, "bad request")
				//return c.Render(code, name, devContext.DevCompParser.Data)
			}

			if requestData.Render != "" && requestData.Component != 0 && requestData.Example != 0 {
				// Parse all templates
				tBuilder := template.Must(template.New("").Funcs(app.DefaultFuncMap).ParseGlob("server/views/*/*.html"))
				// tBuilder = template.Must(tBuilder.ParseGlob("server/views/*/*/*.html"))

				tNew := &echo_custom.Template{
					Templates: tBuilder,
				}

				dataFromJson := devContext.DevCompParser.Data
				requestedCategoryData,err1 := utils.SafeIndexAccess[Category](dataFromJson,requestData.Category)
				requestedComponentData,err2 := utils.SafeIndexAccess[Components](requestedCategoryData.Components, requestData.Component-1)
				requestedExampleData,err3 := utils.SafeIndexAccess[Examples](requestedComponentData.Examples,requestData.Example-1)

				if(err1 || err2 || err3){
					return echo.NewHTTPError(http.StatusNotFound, "Page not found")
				}

				var contentBuffer bytes.Buffer
				tNew.Templates.ExecuteTemplate(&contentBuffer, requestData.Render, requestedExampleData.Data)

				newData := map[string]interface{}{
					"Data":           devContext.DevCompParser.Data,
					"Content":        template.HTML(contentBuffer.String()),
					"Category":       requestedCategoryData,
					"CategoryIndex":  requestData.Category,
					"Component":      requestedComponentData,
					"ComponentIndex": requestData.Component,
					"Example":        requestedExampleData,
					"ExampleIndex":   requestData.Example,
					"ExampleCount":   len(requestedComponentData.Examples),
				}

				return c.Render(code, name+"-render", newData)
			}

			newData := map[string]interface{}{
				"Data": devContext.DevCompParser.Data,
				"CategoryIndex":  -1,
				"ComponentIndex": -1,
			}

			return c.Render(code, name, newData)
		}
	}

	return func(code int, name string) error {
		return c.Render(http.StatusOK, name, []Category{})
	}
}

/*
Creates the echo.Handler that shall be used on the refresh
route. Whenever someone accesses the route, it will
refresh the JSON, and then the person will be redirected
back to whichever URL it originally came from.
*/
func SetDevComponentsRefreshRoute[T any](app *app.App[T]) echo.HandlerFunc {
	handler := func(c echo.Context) error {
		referer := c.Request().Referer()
		parser, ok := app.Plugins["DevComponentParser"]
		if ok {
			parserConverted, conversionOk := parser.(*DevComponentParser)
			if conversionOk {
				parserConverted.ParseJsonConfigFile()
			}
		}
		return c.Redirect(http.StatusFound, referer)
	}

	return handler
}
