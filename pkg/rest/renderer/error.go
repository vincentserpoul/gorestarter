package renderer

import (
	"context"
	"net/http"

	"github.com/vincentserpoul/gorestarter/pkg/rest/mid"

	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render rendering the error
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	// Putting the request error in context, so it can be picked up by logging
	*r = *r.WithContext(context.WithValue(r.Context(), mid.ErrRequestContextKey, e.Err))
	return nil
}

// ErrInvalidRequest when supplied data is not correct
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
	}
}

// ErrRender when there is a server side issue
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Error rendering response.",
	}
}

// ErrNotFound is the wrapped error for not found resources
var ErrNotFound = &ErrResponse{
	HTTPStatusCode: http.StatusNotFound,
	StatusText:     "Resource not found.",
}
