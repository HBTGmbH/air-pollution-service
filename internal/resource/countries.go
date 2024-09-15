package resource

import (
	"air-pollution-service/internal/store"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type CountryResource struct {
	Storage store.Storage
}

func (rs CountryResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})

	return r
}

// List returns all countries
// @Summary List all available countries
// @Description Returns all countries available in the database
// @Tags country
// @Produce json
// @Router /countries/ [get]
// @Success 200 {object} []CountryResponse
// @Failure 400 {object} ErrResponse
// @Failure 404 {object} ErrResponse
func (rs CountryResource) List(w http.ResponseWriter, r *http.Request) {
	countries := rs.Storage.GetCountries()
	if countries == nil {
		_ = render.Render(w, r, ErrNotFound())
		return
	}

	var response []render.Renderer
	for _, article := range countries {
		response = append(response, newCountryResponse(article))
	}

	_ = render.RenderList(w, r, response)
}

// Get returns a single country
// @Summary Get country by its name
// @Description Returns a single country by name
// @Tags country
// @Produce json
// @Router /countries/{name} [get]
// @Param name path string true "name of the country"
// @Success 200 {object} CountryResponse
// @Failure 400 {object} ErrResponse
// @Failure 404 {object} ErrResponse
func (rs CountryResource) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New(fmt.Sprintf("Country name missing"))))
		return
	}

	country := rs.Storage.GetCountry(name)
	if country == nil {
		_ = render.Render(w, r, ErrNotFound())
		return
	}

	response := newCountryResponse(country)
	_ = render.Render(w, r, response)
}
