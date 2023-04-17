package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// readIDParam gets the id query parameter value of a request.
//
// Developer note:
// This function doesn't need to be a receiver of the "application" type since it does not
// depend on any struct fields. However, it helps with creating a mental model of what
// is in the boundary of the "application" type.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We can use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" parameter from
	// the slice. In our project all movies will have a unique positive integer ID, but
	// the value returned by ByName() is always a string. So we try to convert it to a
	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
	// or is less than 1, we know the ID is invalid so return an error.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id parameter: %w", err)
	}

	return id, nil
}
