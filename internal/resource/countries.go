package resource

import (
	"air-pollution-service/internal/model"
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

// CountryResponse TODO
type countryResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
} // @name CountryResponse

func newCountryResponse(country *model.Country) countryResponse {
	return countryResponse{
		Name: country.Name,
		Code: country.Code,
	}
}

func (hr countryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rs CountryResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", rs.Get)
	})

	return r
}

// List TODO
// @Summary TODO
// @Description TODO
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

// Get TODO
// @Summary TODO
// @Description TODO
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
