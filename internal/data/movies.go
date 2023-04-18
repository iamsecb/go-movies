package data

import (
	"time"

	"github.com/haani-niyaz/go-movies/internal/validator"
)

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
	// Use the Runtime type instead of int32. Note that the "omitempty" directive will
	// still work on this: if the Runtime field has the underlying value 0, then it will
	// be considered empty and omitted -- and the MarshalJSON() method we just made
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"` // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

// Validate runs validations on the movie data type.
func (m *Movie) Validate(v *validator.Validator) {
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Note that we're using the Unique() helper in the line below to check that all
	// values in the movies.Genres slice are unique.
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")
}
