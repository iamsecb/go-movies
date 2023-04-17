package data

import "time"

// Annotate the Movie struct with struct tags to control how the keys appear in the
// JSON-encoded output.
//
// Developer note:
// "-" json tag can be used to never present the field in the JSON output to hide
// internal or sensitive information.
//
// "omitempty" json tag will hide the field if the field value is empty. Empty is defined by:
// * Equal to false, 0, or ""
// * An empty array, slice or map
// * A nil pointer or a nil interface value
//
// If you want to use the struct field name as-is and use "omitempty", a leading "," is still needed.
// e.g: "json:,omitempty".
//
// You can force a field to be string for int*, uint*, float* or bool types by doing:
// "`json:runtime,string`".
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"` // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}
