package resource

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"status"`
	StatusText     string `json:"text"`
}

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
