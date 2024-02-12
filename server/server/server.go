package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/UPSxACE/go-local-diary/server/modules/echo_custom"
	"github.com/UPSxACE/go-local-diary/server/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/* Map that will hold the plugins data and state */
type PluginsData = map[string]interface{}

type Server struct {
	t        echo.Renderer
	DevMode  bool
	Echo     *echo.Echo
	Database *db_sqlite3.Database_Sqlite3
	Services *services.Services
	Plugins  PluginsData
}

func (server *Server) Init() {
	// Load plugins
	if server.DevMode {
		dev_component_parser.LoadPlugin(server.Echo, server.Plugins)
	}

	// Print server config
	pluginList := make([]string, 0, len(server.Plugins))
	for pluginName := range server.Plugins {
		pluginList = append(pluginList, pluginName)
	}

	fmt.Println("App Config:")
	if server.DevMode {
		fmt.Println("Dev Mode Enabled")
	}
	fmt.Printf("Extra Plugins: %v\n", pluginList)

	// Start server
	if server.DevMode {
		server.Echo.Use(middleware.Logger())
	}
	server.Echo.Use(middleware.Recover())
	fmt.Println("Initializing echo HTTP server...")
	server.Echo.Logger.Fatal(server.Echo.Start(":1323"))
}

func NewServer(database *db_sqlite3.Database_Sqlite3, devmode bool) (server *Server) {
	server = &Server{
		DevMode: devmode,
		Plugins: map[string]any{},
		Database: database,
	}
	server.Services = services.NewServices(server.Database)

	// Create echo instance
	server.Echo = echo.New()

	// Set the custom context
	server.Echo.Use(echo_custom.GenerateCustomContextMiddleware(server.Services))

	// Setup the normal app configs
	server.setupRenderer()
	server.setupConfig()

	// Prepare the database
	sqlFileReader, err := db_sqlite3.OpenSqlFile("./server/sql/initial.sql")
	if err != nil {
		log.Fatal(err)
	}
	queryThatFailed, err := sqlFileReader.ExecuteAll(server.Database)
	if err != nil {
		log.Fatal(err, queryThatFailed)
	}

	fmt.Printf("Database Tables: %v\n", server.Database.GetTables())

	// FIXME server = &Server{Echo: echoInstance, Database: databaseInstance, Plugins: plugins}
	// Routes
	server.setRoutes()

	return server
}

func (server *Server) setRoutes() {
	server.setIndexRoutes()
	server.setNoteRoutes()
	server.setDevRoutes()
}

// Middleware used in developer mode so the js and css files aren't cached.
func preventCacheMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		c.Response().Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		c.Response().Header().Set("Expires", "0")                                         // Proxies.
		return next(c)
	}
}

// Error handling pages.
func (server *Server) customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	if code == 404 {
		c.Redirect(http.StatusFound, "/404")
	} else {
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}
}

// Setups the template renderer and attaches it to the echo instance.
func (server *Server) setupRenderer() {
	var t echo.Renderer

	if server.DevMode {
		t = &echo_custom.TemplateDevMode{}
	}
	if !server.DevMode {
		// Pre-compile templates in views subdirectories, and subdirectories of those subdirectories
		tBuilder := template.Must(template.New("").Funcs(app.DefaultFuncMap).ParseGlob("server/views/*/*.html"))
		// tBuilder = template.Must(tBuilder.ParseGlob("server/views/*/*/*.html"))
		t = &echo_custom.Template{
			Templates: tBuilder,
		}
	}

	server.t = t
}

// Setups CORS, the middlewares, and the route /public to serve static files.
func (server *Server) setupConfig() {
	server.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323"},
		AllowHeaders: []string{"*"},
		// echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept
	}))

	// Set custom error handler (error pages)
	server.Echo.HTTPErrorHandler = server.customHTTPErrorHandler

	// Middleware
	if server.DevMode {
		server.Echo.Use(preventCacheMiddleware)
	}

	// Serve static files
	server.Echo.Static("/public", "server/public")

	server.Echo.Renderer = server.t
}
