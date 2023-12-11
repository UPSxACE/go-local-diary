package server

import (
	"html/template"

	"github.com/UPSxACE/go-local-diary/server/controllers"
	"github.com/UPSxACE/go-local-diary/template_renderer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	// Pre-compile templates in views subdirectories, and subdirectories of those subdirectories
	tBuilder := template.Must(template.ParseGlob("server/views/*/*.html"))
	// tBuilder = template.Must(tBuilder.ParseGlob("server/views/*/*/*.html"))
	t := &template_renderer.Template{
		Templates: tBuilder,
	}

	// Echo instance
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323"},
		AllowHeaders: []string{"*"},
		// echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept
	}))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve static files
	e.Static("/public", "server/public")

	e.Renderer = t

	// Routes
	controllers.SetIndexRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}