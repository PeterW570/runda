package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) healthCheck(c echo.Context) error {
	data := map[string]string{
		"Status": "OK",
	}

	app.logger.Info("Health check")

	return c.JSON(http.StatusOK, data)
}
