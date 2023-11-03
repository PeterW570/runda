package main

import (
	"errors"
	"strconv"
	"strings"

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

func (app *application) readCSV(c echo.Context, key string, defaultValue []string) []string {
	csv := c.QueryParam(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(c echo.Context, key string, defaultValue int) (int, error) {
	s := c.QueryParam(key)

	if s == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue, err
	}

	return i, nil
}
