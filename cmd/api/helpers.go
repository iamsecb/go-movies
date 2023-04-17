package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Define an envelope type to wrap JSON messages if needed. This is useful when returning errors
// with custom messages.
type envelope map[string]any

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

var (
	errBadJSON           = errors.New("body contains badly-formed JSON")
	errEmptyBody         = errors.New("body must not be empty")
	errTooManyJSONValues = errors.New("body must only contain a single JSON value")
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		// Start triage
		var (
			syntaxError *json.SyntaxError
			// A JSON value is not appropriate for the destination Go type.
			unmarshalTypeError *json.UnmarshalTypeError
			// The decode destination is not valid (usually because it is not a pointer).
			invalidUnmarshalError *json.InvalidUnmarshalError
		)

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			//nolint:goerr113 // Dynamic error is needed
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

			// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
			// for syntax errors in the JSON. So we check for this using errors.Is() and
			// return a generic error message. There is an open issue regarding this at
			// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errBadJSON

			// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
			// JSON value is the wrong type for the target destination. If the error relates
			// to a specific field, then we include that in our error message to make it
			// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				//nolint:goerr113 // Dynamic error is needed
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			//nolint:goerr113 // Dynamic error is needed
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errEmptyBody

			// If the JSON contains a field which cannot be mapped to the target destination
			// then Decode() will now return an error message in the format "json: unknown
			// field "<name>"". We check for this, extract the field name from the error,
			// and interpolate it into our custom error message. Note that there's an open
			// issue at https://github.com/golang/go/issues/29035 regarding turning this
			// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			//nolint:goerr113 // Dynamic error is needed
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic,
		// rather than returning an error to our handler. At the end of this chapter
		// we'll talk about panicking versus returning errors, and discuss why it's an
		// appropriate thing to do in this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			//nolint:wrapcheck
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if errors.Is(err, io.EOF) {
		return errTooManyJSONValues
	}

	return nil
}

// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("writeJSON: an error occurred marshaling json: %w", err)
	}

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for k, v := range headers {
		w.Header()[k] = v
	}

	// Set the "Content-Type: application/json" header on the response. If you forget to
	// this, Go will default to sending a "Content-Type: text/plain; charset=utf-8"
	// header instead.
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	// Write the JSON as the HTTP response body.
	// See following link for why it might be beneficial to log the error.
	// https://stackoverflow.com/a/43976633/2180697
	if _, err := w.Write(js); err != nil {
		app.logger.Println("an error occurred writing response: %w", err)
	}

	return nil
}
