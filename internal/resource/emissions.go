package resource

import (
	"air-pollution-service/internal/model"
	"air-pollution-service/internal/store"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
)

type EmissionResource struct {
	Storage store.Storage
}

func (rs EmissionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/year/", func(r chi.Router) {
		r.Get("/", rs.ListByYear)
		r.Get("/{year}", rs.GetByYear)
	})

	r.Route("/country/", func(r chi.Router) {
		r.Get("/", rs.ListByCountry)
		r.Get("/{id}", rs.GetByCountry)
	})

	return r
}

// ListByYear returns the emissions of all countries accumulated for each year
// @Summary List emissions of each year accumulated over all countries
// @Description All historical emissions data of each year accumulated over all countries
// @Tags emission year
// @Produce json
// @Router /emissions/year/ [get]
// @Success 200 {object} AirPollutionYears
func (rs EmissionResource) ListByYear(w http.ResponseWriter, r *http.Request) {
	response := AirPollutionYears{
		Years: make(map[int]AirPollutionResponse),
	}

	for year, emissions := range rs.Storage.FindAllByYears() {
		response.Years[year] = newAirPollutionEmissionsResponse(emissions)
	}

	_ = render.Render(w, r, response)
}

// GetByYear returns the emissions of all countries accumulated for a single year
// @Summary Get emissions of all countries accumulated for a single year
// @Description All historical emissions data of a year accumulated over all countries, available in the database
// @Tags emission year
// @Produce json
// @Router /emissions/year/{year} [get]
// @Param year path string true "year"
// @Success 200 {object} AirPollutionResponse
// @Failure 400 {object} ErrResponse
func (rs EmissionResource) GetByYear(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New(fmt.Sprintf("Year missing"))))
		return
	}

	var yearEmissions []*model.Emissions
	for _, emissions := range rs.Storage.FindAllByYear(year) {
		yearEmissions = append(yearEmissions, emissions)
	}

	response := newAirPollutionEmissionsResponse(yearEmissions)
	_ = render.Render(w, r, response)
}

// ListByCountry returns the emissions of all years accumulated for each country
// @Summary List emissions of each country accumulated over all years
// @Description All historical emissions data of each country accumulated over all years
// @Tags emission country
// @Produce json
// @Router /emissions/country/ [get]
// @Success 200 {object} AirPollutionCountries
func (rs EmissionResource) ListByCountry(w http.ResponseWriter, r *http.Request) {
	response := AirPollutionCountries{
		Countries: make(map[string]AirPollutionResponse),
	}

	for country, emissions := range rs.Storage.FindAllByCountries() {
		response.Countries[country] = newAirPollutionEmissionsResponse(emissions)
	}

	_ = render.Render(w, r, response)
}

// GetByCountry returns the emissions of all years accumulated for a single country
// @Summary Get emissions of all years accumulated for a single country
// @Description All historical emissions data of a country accumulated over all years, available in the database
// @Tags emission country
// @Produce json
// @Router /emissions/country/{id} [get]
// @Param id path string true "id of the country, case-insensitive"
// @Success 200 {object} AirPollutionResponse
// @Failure 400 {object} ErrResponse
// @Failure 404 {object} ErrResponse
func (rs EmissionResource) GetByCountry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New(fmt.Sprintf("Country id missing"))))
		return
	}

	var countryEmissions []*model.Emissions
	for _, emissions := range rs.Storage.FindAllByCountry(strings.ToLower(id)) {
		countryEmissions = append(countryEmissions, emissions)
	}

	response := newAirPollutionEmissionsResponse(countryEmissions)
	_ = render.Render(w, r, response)
}
