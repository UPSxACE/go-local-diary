package server

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/plugins/dev_component_parser"
	"github.com/labstack/echo/v4"
)

func (server *Server) setDevRoutes() {
	// TODO: migrate this code to the plugin package,
	// and make the ServerApp abstraction to dinamically read this type of functions
	server.Echo.GET("/dev", server.getDev)
	server.Echo.GET("/dev/components", dev_component_parser.SetDevControllerWrapper(server.getDevComponents, server.Plugins))
	server.Echo.GET("/dev/components/refresh", dev_component_parser.SetDevComponentsRefreshRoute(server.Plugins))
}

func (server *Server) getDev(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/dev/components")
}

func (server *Server) getDevComponents(c echo.Context) error {
		renderFunc := dev_component_parser.GetDevComponentParserRenderFunc(c)
		return renderFunc(http.StatusOK, "dev-components")
}