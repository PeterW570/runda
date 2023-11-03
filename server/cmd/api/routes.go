package main

import (
	"github.com/labstack/echo/v4"
)

func (app *application) addRoutes(e *echo.Echo) {
	e.GET("/v1/status", app.healthCheck)

	e.GET("/v1/courses", app.listCourses)
	e.GET("/v1/courses/:id", app.getCourse)
	e.POST("/v1/courses", app.createCourse)
	e.PATCH("/v1/courses/:id", app.updateCourse)
	e.DELETE("/v1/courses/:id", app.deleteCourse)
}
