package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Course struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Version       int32     `json:"version"`
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description,omitempty"`
	Location      Coords    `json:"location" validate:"required"`
	Tags          []string  `json:"tags"`
	Website       string    `json:"website,omitempty" validate:"optional_uri"`
}

type CourseModel struct {
	DB *sqlx.DB
}

func (c CourseModel) Insert(course *Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
        INSERT INTO courses (name, description, location, tags, website) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, last_updated_at, version`

	args := []any{
		course.Name,
		course.Description,
		course.Location.AsPostgresPointString(),
		pq.Array(course.Tags),
		course.Website,
	}

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&course.ID, &course.CreatedAt, &course.LastUpdatedAt, &course.Version)
}

func (c CourseModel) Get(id int64) (*Course, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
        SELECT id, created_at, last_updated_at, version, name, description, location[0] as longitude, location[1] as latitude, tags, website
        FROM courses
        WHERE id = $1`

	var course Course

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.CreatedAt,
		&course.LastUpdatedAt,
		&course.Version,
		&course.Name,
		&course.Description,
		&course.Location.Longitude,
		&course.Location.Latitude,
		pq.Array(&course.Tags),
		&course.Website,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &course, nil
}

func (c CourseModel) Update(course *Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
        UPDATE courses 
        SET name = $1, description = $2, location = $3, tags = $4, website = $5, last_updated_at = now(), version = version + 1
        WHERE id = $6 AND version = $7
        RETURNING version, last_updated_at`

	args := []any{
		course.Name,
		course.Description,
		course.Location.AsPostgresPointString(),
		pq.Array(course.Tags),
		course.Website,
		course.ID,
		course.Version,
	}

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&course.Version, &course.LastUpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (c CourseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
        DELETE FROM courses
        WHERE id = $1`

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
