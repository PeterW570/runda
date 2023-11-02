package database

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Courses CourseModel
}

func NewModels(db *DB) Models {
	return Models{
		Courses: CourseModel{DB: db.DB},
	}
}
