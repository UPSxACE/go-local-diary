package controllers

import (
	"net/http"

	"github.com/UPSxACE/go-local-diary/server/internal/models"
	"github.com/UPSxACE/go-local-diary/server/pkg/echo_custom"
	"github.com/labstack/echo/v4"
)

func SetIndexRoutes(e *echo.Echo) {
	e.GET("/", echo_custom.RedirectNotConfiguredToWelcomeMiddleware(GetIndexController))
	e.GET("/welcome", echo_custom.RedirectNotConfiguredToWelcomeMiddleware(GetWelcomeController))
	e.POST("/welcome", echo_custom.RedirectNotConfiguredToWelcomeMiddleware(PostWelcomeController))
	e.GET("/404", Get404Controller)
}

func GetIndexController(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func GetWelcomeController(c echo.Context) error {
	cc := c.(*echo_custom.CustomEchoContext)

	if cc.IsConfigured {
		return cc.Redirect(http.StatusMovedPermanently, "/")
	}

	return cc.Render(http.StatusOK, "welcome", nil)
}

func PostWelcomeController(c echo.Context) error {
	cc := c.(*echo_custom.CustomEchoContext)

	if cc.IsConfigured {
		return cc.Redirect(http.StatusMovedPermanently, "/")
	}

	store, err := models.CreateStoreAppConfig(cc.App, true, cc.Request().Context())
	if err != nil {
		return err;
	}

	exists, model, err := store.CheckByName("configured")
	if err != nil {
		return err
	}

	if !exists {
		model = models.AppConfigModel{Name: "configured", Value: "1"}
		model, err = store.Create(model)
		if err != nil {
			return err
		}
	}
	if exists {
		model.Value = "1";
		model, err = store.UpdateByName(model.Name, model)
		if err != nil{
			return err;
		}
	}

	errs := store.Close()
	if(len(errs) != 0){
		return err;
	}

	return cc.Redirect(http.StatusMovedPermanently, "/")
}

func Get404Controller(c echo.Context) error {
	return c.Render(http.StatusOK, "404", nil)
}
