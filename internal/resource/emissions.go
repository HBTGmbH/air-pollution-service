package resource

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/go-chi/chi/v5"
)

// EmissionResource serves statistics about the emissions of the data set.
type EmissionResource struct {
	Storage Storage
}

// Routes returns the router of all emission endpoints.
func (rs EmissionResource) Routes() *chi.Mux {
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

// ListByYear returns the emissions of all countries accumulated for each year.
//
//	@Summary		List emissions of each year accumulated over all countries
//	@Description	All historical emissions data of each year accumulated over all countries
//	@Tags			emission year
//	@Produce		json
//	@Router			/emissions/year/ [get]
//	@Success		200	{object}	AirPollutionYears
func (rs EmissionResource) ListByYear(w http.ResponseWriter, r *http.Request) {
	response := AirPollutionYears{
		Years: make(map[int]AirPollutionResponse),
	}

	for year, emissions := range rs.Storage.FindAllByYears() {
		yearResponse, err := newAirPollutionResponse(emissions)
		if err != nil {
			respondInternalError(w, r, fmt.Errorf("year %d: %w", year, err))

			return
		}
		response.Years[year] = yearResponse
	}

	respond(w, r, response)
}

// GetByYear returns the emissions of all countries accumulated for a single year.
//
//	@Summary		Get emissions of all countries accumulated for a single year
//	@Description	All historical emissions data of a year accumulated over all countries, available in the database
//	@Tags			emission year
//	@Produce		json
//	@Router			/emissions/year/{year} [get]
//	@Param			year	path		string	true	"year"
//	@Success		200		{object}	AirPollutionResponse
//	@Failure		400		{object}	ErrResponse
//	@Failure		404		{object}	ErrResponse
func (rs EmissionResource) GetByYear(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(chi.URLParam(r, "year"))
	if err != nil {
		respond(w, r, ErrInvalidRequest("Year missing"))

		return
	}

	respondWithAirPollution(w, r, slices.Collect(maps.Values(rs.Storage.FindAllByYear(year))))
}

// ListByCountry returns the emissions of all years accumulated for each country.
//
//	@Summary		List emissions of each country accumulated over all years
//	@Description	All historical emissions data of each country accumulated over all years
//	@Tags			emission country
//	@Produce		json
//	@Router			/emissions/country/ [get]
//	@Success		200	{object}	AirPollutionCountries
func (rs EmissionResource) ListByCountry(w http.ResponseWriter, r *http.Request) {
	response := AirPollutionCountries{
		Countries: make(map[string]AirPollutionResponse),
	}

	for country, emissions := range rs.Storage.FindAllByCountries() {
		countryResponse, err := newAirPollutionResponse(emissions)
		if err != nil {
			respondInternalError(w, r, fmt.Errorf("country %s: %w", country, err))

			return
		}
		response.Countries[country] = countryResponse
	}

	respond(w, r, response)
}

// GetByCountry returns the emissions of all years accumulated for a single country.
//
//	@Summary		Get emissions of all years accumulated for a single country
//	@Description	All historical emissions data of a country accumulated over all years, available in the database
//	@Tags			emission country
//	@Produce		json
//	@Router			/emissions/country/{id} [get]
//	@Param			id	path		string	true	"id of the country, case-insensitive"
//	@Success		200	{object}	AirPollutionResponse
//	@Failure		400	{object}	ErrResponse
//	@Failure		404	{object}	ErrResponse
func (rs EmissionResource) GetByCountry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respond(w, r, ErrInvalidRequest("Country id missing"))

		return
	}

	emissions := rs.Storage.FindAllByCountry(strings.ToLower(id))
	respondWithAirPollution(w, r, slices.Collect(maps.Values(emissions)))
}

// respondWithAirPollution responds with the statistics over all given
// emissions, or with 404 if there are none.
func respondWithAirPollution(w http.ResponseWriter, r *http.Request, emissions []*model.Emissions) {
	if len(emissions) == 0 {
		respond(w, r, ErrNotFound())

		return
	}

	response, err := newAirPollutionResponse(emissions)
	if err != nil {
		respondInternalError(w, r, err)

		return
	}

	respond(w, r, &response)
}
