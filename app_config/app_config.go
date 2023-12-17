/*
The app_config package holds the interfaces and structs that will hold

the application configuration state, instances and plugins,and will be shared
across the entire app.
*/
package app_config

import (
	"html/template"
	"strconv"
	"strings"
)

/*
Holds the app configuration state. One instance per app.

The app must be initialized with a Database, whichever type it is.
*/
type AppConfig[DatabaseGeneric any] struct {
	Database *DatabaseGeneric
	DevMode  bool
	Plugins  PluginsData
}

/* Map that will hold the plugins data and state */
type PluginsData = map[string]interface{}

/* Some useful extra funcs for the templating system of html/template. */
var DefaultFuncMap template.FuncMap = template.FuncMap{
	"list":       list,
	"obj":        obj,
	"sum":        sum,
	"sumStr":     sumStr,
	"htmlbreaks": htmlBreaks,
}

type DefMapInvalidArgs struct{
	code int;
	message string;
}

func (m *DefMapInvalidArgs) Error() string {
	return m.message;
}

/* Returns everything it got as argument, inside a slice. */
func list(args ...interface{}) []interface{} {
	slice := []interface{}{}
	slice = append(slice, args...)
	return slice
}

/*
Transforms a string into a map.

ej: "key1:val,,key2: this is a sentence -> {key1: "val", key2: "this is a sentence"}
*/
func obj(str string) map[string]string {
	obj := map[string]string{}

	pairs := strings.Split(str, ",,")
	for _, pair := range pairs {
		keyVal := strings.Split(pair, ":")

		if(len(keyVal) != 2){
			panic(&DefMapInvalidArgs{1, "The strings used as arguments are not well formated"})
		}

		key := keyVal[0]
		val := keyVal[1]

		obj[key] = val
	}

	return obj
}

/* Sums two integers. */
func sum(num1 int, num2 int) int {
	return num1 + num2
}

/* Transforms two strings into integers, then sums. */
func sumStr(num1 string, num2 string) int {
	num1Converted, err := strconv.Atoi(num1)
	num2Converted, err2 := strconv.Atoi(num2)
	if err != nil || err2 != nil {
		panic(&DefMapInvalidArgs{2, "Couldn't convert argument to integer"})
	}
	return num1Converted + num2Converted
}

/*
Transforms a string into html, but ONLY parses \n to <br>.
All other html tags will be ignored.
*/
func htmlBreaks(str string) template.HTML {
	safeHtml := template.HTMLEscapeString(str)
	safeHtmlWithBreaks := strings.ReplaceAll(safeHtml, "\n", "<br>")
	finalHtml := template.HTML(safeHtmlWithBreaks)
	return finalHtml
}
