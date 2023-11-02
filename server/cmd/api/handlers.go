package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) healthCheck(c echo.Context) error {
	data := envelope{
		"Status":  "OK",
		"Version": version,
	}

	return c.JSON(http.StatusOK, data)
}
