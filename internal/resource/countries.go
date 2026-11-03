// Package resource implements the HTTP handlers of the REST API.
package resource

import (
	"net/http"
	"strings"

	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Storage provides read access to the data set.
type Storage interface {
	FindAllByYears() map[int][]*model.Emissions
	FindAllByYear(year int) map[string]*model.Emissions
	FindAllByCountries() map[string][]*model.Emissions
	FindAllByCountry(id string) map[int]*model.Emissions
	GetCountry(id string) *model.Country
	GetCountries() []*model.Country
}

// CountryResource serves the countries of the data set.
type CountryResource struct {
	Storage Storage
}

// Routes returns the router of all country endpoints.
func (rs CountryResource) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})

	return r
}

// List returns all countries.
//
//	@Summary		List all available countries
//	@Description	Returns all countries available in the database
//	@Tags			country
//	@Produce		json
//	@Router			/countries/ [get]
//	@Success		200	{object}	[]CountryResponse
//	@Failure		400	{object}	ErrResponse
//	@Failure		404	{object}	ErrResponse
func (rs CountryResource) List(w http.ResponseWriter, r *http.Request) {
	countries := rs.Storage.GetCountries()
	if len(countries) == 0 {
		respond(w, r, ErrNotFound())

		return
	}

	response := make([]render.Renderer, 0, len(countries))
	for _, country := range countries {
		response = append(response, newCountryResponse(country))
	}

	respondList(w, r, response)
}

// Get returns a single country.
//
//	@Summary		Get country by its id
//	@Description	Returns a single country by id
//	@Tags			country
//	@Produce		json
//	@Router			/countries/{id} [get]
//	@Param			id	path		string	true	"ID of the country, case-insensitive"
//	@Success		200	{object}	CountryResponse
//	@Failure		400	{object}	ErrResponse
//	@Failure		404	{object}	ErrResponse
func (rs CountryResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respond(w, r, ErrInvalidRequest("Country id missing"))

		return
	}

	country := rs.Storage.GetCountry(strings.ToLower(id))
	if country == nil {
		respond(w, r, ErrNotFound())

		return
	}

	respond(w, r, newCountryResponse(country))
}
