package resource

import (
	"github.com/go-chi/render"
	"net/http"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	HTTPStatusCode int    `json:"status"`
	StatusText     string `json:"text"`
} // @name ErrResponse

func ErrRender(message string, status int) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: status,
		StatusText:     message,
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
