package server

import (
	"html/template"
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/controllers"
	"github.com/UPSxACE/go-local-diary/server/template_renderer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(appInstance *app.App[db_sqlite3.Database_Sqlite3]) {
	t := setupRenderer(appInstance);

	// Create echo instance
	e := echo.New()

	setupConfig(appInstance, e, &t)	

	// Routes
	controllers.SetIndexRoutes(e)
	if appInstance.DevMode {
		controllers.SetDevRoutes(e, appInstance)
	}

	// Start server
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":1323"))
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
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
    }
    c.Logger().Error(err)
	if(code == 404){
		c.Redirect(http.StatusPermanentRedirect, "/404")
	} else {
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}	
}

// Setups the template renderer and attaches it to the echo instance.
func setupRenderer(appInstance *app.App[db_sqlite3.Database_Sqlite3]) echo.Renderer{
	var t echo.Renderer

	if appInstance.DevMode {
		t = &template_renderer.TemplateDevMode{}
	}
	if !appInstance.DevMode {
		// Pre-compile templates in views subdirectories, and subdirectories of those subdirectories
		tBuilder := template.Must(template.New("").Funcs(app.DefaultFuncMap).ParseGlob("server/views/*/*.html"))
		// tBuilder = template.Must(tBuilder.ParseGlob("server/views/*/*/*.html"))
		t = &template_renderer.Template{
			Templates: tBuilder,
		}
	}

	return t;
}

// Setups CORS, the middlewares, and the route /public to serve static files.
func setupConfig(appInstance *app.App[db_sqlite3.Database_Sqlite3], e *echo.Echo, t *echo.Renderer) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323"},
		AllowHeaders: []string{"*"},
		// echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept
	}))
	
	// Set custom error handler (error pages)
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Middleware
	if appInstance.DevMode {
		e.Use(preventCacheMiddleware)
	}

	// Serve static files
	e.Static("/public", "server/public")

	e.Renderer = *t
}

