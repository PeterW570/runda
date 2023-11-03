package validation

import (
	"errors"

	"peterweightman.com/runda/internal/database"
)

var (
	ErrPageBelowMinimum     = errors.New("page below minimum")
	ErrPageAboveMaximum     = errors.New("page above maximum")
	ErrPageSizeBelowMinimum = errors.New("page size below minimum")
	ErrPageSizeAboveMaximum = errors.New("page size above maximum")
	ErrSortInvalid          = errors.New("sort invalid")
)

var (
	MinPage     = 1
	MaxPage     = 10_000_000 - 1
	MinPageSize = 1
	MaxPageSize = 100 - 1
)

func ValidateFilters(f database.Filters) error {
	if f.Page < MinPage {
		return ErrPageBelowMinimum
	}
	if f.Page > MaxPage {
		return ErrPageAboveMaximum
	}
	if f.PageSize < MinPageSize {
		return ErrPageSizeBelowMinimum
	}
	if f.PageSize > MaxPageSize {
		return ErrPageSizeAboveMaximum
	}

	if !In(f.Sort, f.SortSafelist...) {
		return ErrSortInvalid
	}

	return nil
}
