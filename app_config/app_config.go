package app_config

import (
	"html/template"
	"strconv"
	"strings"
)

type AppConfig[DatabaseGeneric any] struct {
	Database *DatabaseGeneric
	DevMode  bool
	Plugins  PluginsData
}

type PluginsData = map[string]interface{}

var DefaultFuncMap template.FuncMap = template.FuncMap{
	"devmode": func() bool { return false },
	"list": func(args ...interface{}) []interface{} {
		slice := []interface{}{}
		slice = append(slice, args...)
		return slice
	},
	"obj": func(str string) map[string]string {
		obj := map[string]string{}

		pairs := strings.Split(str, ",,")
		for _, pair := range pairs {
			keyVal := strings.Split(pair, ":")

			key := keyVal[0]
			val := keyVal[1]

			obj[key] = val
		}

		return obj
	},
	"sum": func(num1 int, num2 int) int {
		return num1 + num2
	},
	"sumStr": func(num1 string, num2 string) int {
		num1Converted, err := strconv.Atoi(num1)
		num2Converted, err2 := strconv.Atoi(num2)
		if err != nil || err2 != nil {
			return 0
		}
		return num1Converted + num2Converted
	},
	"htmlbreaks": func(str string) template.HTML {
		safeHtml := template.HTMLEscapeString(str)
		safeHtmlWithBreaks := strings.ReplaceAll(safeHtml, "\n", "<br>")
		finalHtml := template.HTML(safeHtmlWithBreaks)
		return finalHtml
	},
}