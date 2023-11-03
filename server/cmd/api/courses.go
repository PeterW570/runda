package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"peterweightman.com/runda/internal/database"
	"peterweightman.com/runda/internal/validation"
)

func (app *application) createCourse(c echo.Context) error {
	course := new(database.Course)
	err := c.Bind(course)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(course); err != nil {
		return err
	}

	err = app.models.Courses.Insert(course)
	if err != nil {
		app.logger.Error("Error inserting course", "error", err)
		return echo.ErrInternalServerError
	}

	c.Response().Header().Set(echo.HeaderLocation, fmt.Sprintf("/v1/courses/%d", course.ID))
	return c.JSON(http.StatusCreated, envelope{"course": course})
}

func (app *application) getCourse(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "course not found")
	}

	course, err := app.models.Courses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return echo.ErrNotFound
		default:
			app.logger.Error("Error getting course", "error", err)
			return echo.ErrInternalServerError
		}
	}

	return c.JSON(http.StatusOK, envelope{"course": course})
}

func (app *application) updateCourse(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		return echo.ErrNotFound
	}

	course, err := app.models.Courses.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return echo.ErrNotFound
		default:
			return echo.ErrInternalServerError
		}
	}

	input := new(database.Course)
	err = c.Bind(input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if input.Name != "" {
		course.Name = input.Name
	}

	if input.Description != "" {
		course.Description = input.Description
	}

	if input.Location.Latitude != 0 || input.Location.Longitude != 0 {
		course.Location = input.Location
	}

	if input.Tags != nil {
		course.Tags = input.Tags
	}

	if input.Website != "" {
		course.Website = input.Website
	}

	if err = c.Validate(course); err != nil {
		return err
	}

	err = app.models.Courses.Update(course)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrEditConflict):
			return echo.NewHTTPError(http.StatusConflict, "unable to update the record due to an edit conflict, please try again")
		default:
			app.logger.Error("Error updating course", "error", err)
			return echo.ErrInternalServerError
		}
	}

	return c.JSON(http.StatusOK, envelope{"course": course})
}

func (app *application) deleteCourse(c echo.Context) error {
	id, err := app.readIDParam(c)
	if err != nil {
		return echo.ErrNotFound
	}

	err = app.models.Courses.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return echo.ErrNotFound
		default:
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

func (app *application) listCourses(c echo.Context) error {
	var input struct {
		Name string
		Tags []string
		database.Filters
	}

	input.Name = c.QueryParam("name")
	input.Tags = app.readCSV(c, "tags", []string{})

	var err error
	input.Filters.Page, err = app.readInt(c, "page", 1)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page number")
	}
	input.Filters.PageSize, err = app.readInt(c, "page_size", 20)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page size")
	}

	input.Filters.Sort = c.QueryParam("sort")
	if input.Filters.Sort == "" {
		input.Filters.Sort = "id"
	}

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if err = validation.ValidateFilters(input.Filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	courses, metadata, err := app.models.Courses.GetAll(input.Name, input.Tags, input.Filters)
	if err != nil {
		app.logger.Error("Error getting courses", "error", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, envelope{"courses": courses, "metadata": metadata})
}
