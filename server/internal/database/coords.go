package database

import "fmt"

type Coords struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

func (c Coords) AsPostgresPointString() string {
	return fmt.Sprintf("(%f, %f)", c.Longitude, c.Latitude)
}
