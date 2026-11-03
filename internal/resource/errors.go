package resource

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	StatusText string `json:"text,omitempty"`
	ErrText    string `example:"The requested resource was not found on the server." json:"error,omitempty"`
	StatusCode int    `json:"-"`
	ErrCode    int    `example:"404"                                                 json:"code,omitempty"`
} //	@name	ErrResponse

// Render implements render.Renderer and sets the HTTP status of the response.
func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)

	return nil
}

// ErrInvalidRequest returns a structured http response in case of an invalid
// request. The message is returned to the client verbatim.
func ErrInvalidRequest(message string) *ErrResponse {
	return &ErrResponse{
		ErrText:    message,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request",
	}
}

// ErrNotFound returns a structured http response in case a resource was not found.
func ErrNotFound() *ErrResponse {
	return &ErrResponse{
		StatusCode: http.StatusNotFound,
		StatusText: "Resource not found",
	}
}

// ErrInternal returns a structured http response in case of an unexpected
// server-side failure.
func ErrInternal() *ErrResponse {
	return &ErrResponse{
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal server error",
	}
}

// respond renders v, or an internal server error if rendering fails.
func respond(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	err := render.Render(w, r, v)
	if err != nil {
		respondInternalError(w, r, fmt.Errorf("rendering response: %w", err))
	}
}

// respondList renders l, or an internal server error if rendering fails.
func respondList(w http.ResponseWriter, r *http.Request, l []render.Renderer) {
	err := render.RenderList(w, r, l)
	if err != nil {
		respondInternalError(w, r, fmt.Errorf("rendering response: %w", err))
	}
}

// respondInternalError logs err and responds with an internal server error
// that does not leak any details of err to the client.
func respondInternalError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Responding with internal server error: %s", err)
	response := ErrInternal()
	render.Status(r, response.StatusCode)
	render.Respond(w, r, response)
}
