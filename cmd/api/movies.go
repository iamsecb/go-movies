package main

import (
	"fmt"
	"net/http"
)

// Add a createMovieHandler for the "POST api/v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// TO-DO: Store movie
	fmt.Fprint(w, "create a new movie")
}

// Add a showMovieHandler for the "GET api/v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// If ID is invalid so we use the http.NotFound() function to return a 404 Not Found response.
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	// TO-DO: Get movie.
	fmt.Fprintf(w, "show details of movie %d\n", id)
}