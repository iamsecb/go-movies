package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/haani-niyaz/go-movies/internal/data"
)

// Add a createMovieHandler for the "POST api/v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	//
	// If malformed JSON is sent, we need to decide how much of it is returned in the response.
	// If this is a private API, returning the errors as-is is probably safe to do so. If it is
	// public on the other hand we should "information hide" to not expose any internal
	// implementation.
	// You can also run into a range of issues like the level of details might be insufficient,
	// inconsistent language etc. so the errors should be triaged.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET api/v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// If ID is invalid so we use the http.NotFound() function to return a 404 Not Found response.
	id, err := app.readIDParam(r)
	if err != nil {
		app.errNotFoundResponse(w, r)

		return
	}

	// TO-DO: Get movie.
	m := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   data.Runtime(102),
		Genres:    []string{"drama", "romance", "war"},
	}

	if err := app.writeJSON(w, http.StatusOK, &m, nil); err != nil {
		app.errServerResponse(w, r, err)
	}
}
