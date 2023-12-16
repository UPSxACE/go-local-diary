// Web
package dev_component_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/UPSxACE/go-local-diary/server/template_renderer"
	"github.com/UPSxACE/go-local-diary/utils"

	"github.com/labstack/echo/v4"
)

type DevComponentParser struct {
	Data []Category
}

func LoadPlugin[T any](appConfig *app_config.AppConfig[T]) {
	fmt.Println("Loading Plugin DevComponentParser...")
	if appConfig.DevMode {
		appConfig.Plugins["DevComponentParser"] = Init()
	}
}

func Init() *DevComponentParser {
	devCompPars := DevComponentParser{}
	devCompPars.ParseJsonConfigFile()
	return &devCompPars
}

func (devComponentParser *DevComponentParser) ParseJsonConfigFile() *DevComponentParser {
	jsonFile, err := os.Open("./server/dev-components.json")
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

func SetDevControllerWrapper[T any](controller echo.HandlerFunc, appConfig *app_config.AppConfig[T]) echo.HandlerFunc {
	parser, ok := appConfig.Plugins["DevComponentParser"]
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

type DevComponentParserRequest struct {
	Render    string `query:"render"`
	Category  int    `query:"ct"` //index
	Component int    `query:"cp"` // index+1
	Example   int    `query:"e"`  // index+1
}

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
				tBuilder := template.Must(template.New("").Funcs(app_config.DefaultFuncMap).ParseGlob("server/views/*/*.html"))
				// tBuilder = template.Must(tBuilder.ParseGlob("server/views/*/*/*.html"))

				tNew := &template_renderer.Template{
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

type devContextWrapper struct {
	echo.Context
	DevCompParser *DevComponentParser
}

func SetDevComponentsRefreshRoute[T any](appConfig *app_config.AppConfig[T]) echo.HandlerFunc {
	handler := func(c echo.Context) error {
		referer := c.Request().Referer()
		parser, ok := appConfig.Plugins["DevComponentParser"]
		if ok {
			parserConverted, conversionOk := parser.(*DevComponentParser)
			if conversionOk {
				parserConverted.ParseJsonConfigFile()
			}
		}
		return c.Redirect(http.StatusTemporaryRedirect, referer)
	}

	return handler
}
