package resource

import (
	"air-pollution-service/internal/model"
	"air-pollution-service/internal/store"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type CountryResource struct {
	Storage store.Storage
}

type countryResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

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

func (rs CountryResource) List(w http.ResponseWriter, r *http.Request) {
	countries := rs.Storage.GetCountries()
	if countries == nil {
		if err := render.Render(w, r, ErrRender(fmt.Sprintf("No country found"), 404)); err != nil {
			render.Status(r, 500)
		}
		return
	}

	var response []render.Renderer
	for _, article := range countries {
		response = append(response, newCountryResponse(article))
	}

	if err := render.RenderList(w, r, response); err != nil {
		if err != nil {
			render.Status(r, 500)
		}
	}
}

func (rs CountryResource) Get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		if err := render.Render(w, r, ErrRender(fmt.Sprintf("Country name missing"), 400)); err != nil {
			render.Status(r, 500)
		}
		return
	}

	country := rs.Storage.GetCountry(name)
	if country == nil {
		if err := render.Render(w, r, ErrRender(fmt.Sprintf("No country with name %s found", name), 404)); err != nil {
			render.Status(r, 500)
		}
		return
	}

	response := newCountryResponse(country)
	if err := render.Render(w, r, response); err != nil {
		render.Status(r, 500)
	}
}
