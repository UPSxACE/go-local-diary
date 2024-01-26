package echo_custom

import (
	"fmt"
	"net/http"

	"github.com/UPSxACE/go-local-diary/app"
	"github.com/UPSxACE/go-local-diary/plugins/db_sqlite3"
	"github.com/UPSxACE/go-local-diary/server/models"
	"github.com/labstack/echo/v4"
)

/*
Struct to extend the echo.Context fields.
*/
type CustomEchoContext struct {
	echo.Context
	IsConfigured bool
	App *app.App[db_sqlite3.Database_Sqlite3]
}

/*
Generates the middleware to set the custom echo context with its variables with custom behavior:

- fills the field 'IsConfigured' of the context with the correct value
*/
func GenerateCustomContextMiddleware(app *app.App[db_sqlite3.Database_Sqlite3]) func(next echo.HandlerFunc) echo.HandlerFunc {
	var middleware = func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomEchoContext{Context: c, App: app}
			context := cc.Request().Context()

			appConfigStore, err := models.CreateStoreAppConfig(cc.App, true, context)
			if err != nil {
				return err;
			}
			failNotConfigured := func() error {
				cc.IsConfigured = false
				appConfigStore.Close()
				return next(cc)
			}
			internalErr := func(error) error {
				cc.IsConfigured = false
				appConfigStore.Close()
				return err;
			}
			
			exists, model, err := appConfigStore.CheckByName("configured")
			if err != nil {
				return internalErr(err)
			}

			if !exists {
				_, err = appConfigStore.Create(models.AppConfigModel{Id: 0, Name: "configured", Value: "0"})
				if err != nil {
					return internalErr(err)
				}
				return failNotConfigured()
			}

			if model.Value == "1" {
				cc.IsConfigured = true
			} else if model.Value == "0" {
				return failNotConfigured()
			} else {
				fmt.Println("Invalid configuration value in config: 'configured'", err)
				return failNotConfigured()
			}

			// It exists and is configured
			appConfigStore.Close()
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
		
		url := cc.Context.Request().URL.String()

		if(!cc.IsConfigured && url != "/welcome"){
			return cc.Redirect(http.StatusFound, "/welcome")
		}

		return controller(cc)
	}	 

	return fnc
}