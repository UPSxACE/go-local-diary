package echo_custom

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/internal/services"
	"github.com/labstack/echo/v4"
)

/*
Struct to extend the echo.Context fields.
*/
type CustomEchoContext struct {
	echo.Context
	IsConfigured bool
	// NOTE: This (below) seems to be a bad practice, so it's removed for now
	// Context must only carry data tied to the specific request
	// App *app.App[db_sqlite3.Database_Sqlite3]
}

/*
Generates the middleware to set the custom echo context with its variables with custom behavior:

- fills the field 'IsConfigured' of the context with the correct value
*/
func GenerateCustomContextMiddleware(app *app.App[db_sqlite3.Database_Sqlite3]) func(next echo.HandlerFunc) echo.HandlerFunc {
	var middleware = func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomEchoContext{Context: c}
			context := cc.Request().Context()
			
			configured, err := services.AppConfig.IsAppConfigured(app,context)
			if err != nil {
				return err;
			}
			cc.IsConfigured = configured;

			return next(cc)
		}
	}

	return middleware
}


/** Middleware at CONTROLLER LEVEL, to redirect requests from other routes to /welcome
 *  when the app is not configured yet.
 */
func RedirectNotConfiguredToWelcomeMiddleware(controller func (c echo.Context) error) (func (c echo.Context) error){
	fnc := func (c echo.Context) error {
		cc := c.(*CustomEchoContext)
		
		path := cc.Context.Request().URL.Path

		if(!cc.IsConfigured && path != "/welcome"){
			isHtmxBoosted := cc.Request().Header.Get("HX-Boosted") != ""

			if(isHtmxBoosted){
				cc.Response().Header().Set("HX-Redirect", "/welcome")
				return cc.Render(http.StatusOK, "", nil)
			}

			return cc.Redirect(http.StatusFound, "/welcome")
		}

		return controller(cc)
	}	 

	return fnc
}