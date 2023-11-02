package main

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type envelope map[string]any

func (app *application) readIDParam(c echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
