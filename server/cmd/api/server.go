package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"peterweightman.com/runda/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// for _, err := range err.(validator.ValidationErrors) {
		//     fmt.Printf("Validation Error: Field '%s' failed the '%s' tag.\n", err.Field(), err.Tag())
		// }
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (app *application) serveHTTP() error {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validation.NewValidator()}

	e.Use(slogecho.New(app.logger))
	e.Use(middleware.Recover())

	app.addRoutes(e)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      e,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", slog.Group("server", "addr", s.Addr), slog.String("env", app.config.env))

	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
