package dev_component_parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/UPSxACE/go-local-diary/app_config"
	"github.com/labstack/echo/v4"
)

type DevComponentParser struct {
	Data []Category
}

func LoadPlugin(appConfig *app_config.AppConfig) {
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

func SetDevControllerWrapper(controller echo.HandlerFunc, appConfig *app_config.AppConfig) echo.HandlerFunc {
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

func GetDevComponentParserRenderFunc(c echo.Context) func(code int, name string) error {
	devContext, ok := c.(*devContextWrapper)

	if ok {
		return func(code int, name string) error {
			return c.Render(code, name, devContext.DevCompParser.Data)
		}
	}

	return func(code int, name string) error {
		return c.Render(http.StatusOK, "dev-components", []Category{})
	}
}

type devContextWrapper struct {
	echo.Context
	DevCompParser *DevComponentParser
}
