package resource

import (
	"github.com/go-chi/render"
	"net/http"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	StatusCode int    `json:"-"`
	StatusText string `json:"text,omitempty"`
	Err        error  `json:"-"`
	ErrCode    int    `json:"code,omitempty" example:"404"`
	ErrText    string `json:"error,omitempty" example:"The requested resource was not found on the server."`
} // @name ErrResponse

// ErrRender returns a structured http response in case of rendering errors
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: http.StatusUnprocessableEntity,
		StatusText: "Error rendering response",
	}
}

// ErrInvalidRequest returns a structured http response in case of an invalid request
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request",
	}
}

// ErrNotFound returns a structured http response in case a resource was not found
func ErrNotFound() render.Renderer {
	return &ErrResponse{
		StatusCode: http.StatusNotFound,
		StatusText: "Resource not found",
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
