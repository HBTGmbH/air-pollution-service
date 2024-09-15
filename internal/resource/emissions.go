package resource

import (
	"air-pollution-service/internal/model"
	"air-pollution-service/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/montanaflynn/stats"
	"net/http"
	"strconv"
)

type EmissionResource struct {
	Storage store.Storage
}

// AirPollutionResponse TODO
type airPollutionResponse struct {
	Average           float64 `json:"average"`
	Median            float64 `json:"median"`
	StandardDeviation float64 `json:"standard_deviation"`
} // @name AirPollutionResponse

func newAirPollutionResponse(emissions []*model.Emissions, f func(emission *model.Emissions) float64) airPollutionResponse {
	total := 0.0
	values := make([]float64, len(emissions))

	for i, emission := range emissions {
		current := f(emission)
		values[i] = current
		total = total + current
	}

	median, _ := stats.Median(values)
	standardDeviation, _ := stats.StandardDeviation(values)

	return airPollutionResponse{
		Average:           total / float64(len(emissions)),
		Median:            median,
		StandardDeviation: standardDeviation,
	}
}

// AirPollutionEmissionsResponse TODO
type airPollutionEmissionsResponse struct {
	NOxEmissions   airPollutionResponse `json:"nox_emissions"`
	SO2Emissions   airPollutionResponse `json:"sulphur_dioxide_emissions"`
	COEmissions    airPollutionResponse `json:"carbon_monoxide_emissions"`
	OCEmissions    airPollutionResponse `json:"organic_carbon_emissions"`
	NMVOCEmissions airPollutionResponse `json:"nmvoc_emissions"`
	BCEmissions    airPollutionResponse `json:"black_carbon_emissions"`
	NH3Emissions   airPollutionResponse `json:"ammonia_emissions"`
} // @name AirPollutionEmissionsResponse

func newAirPollutionEmissionsResponse(emissions []*model.Emissions) airPollutionEmissionsResponse {
	return airPollutionEmissionsResponse{
		NOxEmissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NOxEmissions
		}),
		SO2Emissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.SO2Emissions
		}),
		COEmissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.COEmissions
		}),
		OCEmissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.OCEmissions
		}),
		NMVOCEmissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NMVOCEmissions
		}),
		BCEmissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.BCEmissions
		}),
		NH3Emissions: newAirPollutionResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NH3Emissions
		}),
	}
}

func (hr airPollutionEmissionsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rs EmissionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/year/", func(r chi.Router) {
		r.Get("/", rs.ListByYear)
		r.Get("/{year}", rs.GetByYear)
	})

	r.Route("/country/", func(r chi.Router) {
		r.Get("/", rs.ListByCountry)
		r.Get("/{name}", rs.GetByCountry)
	})

	return r
}

// ListByYear TODO
// @Summary TODO
// @Description TODO
// @Tags emission year
// @Produce json
// @Router /emissions/year/ [get]
// @Success 200 {object} map[country]AirPollutionEmissionsResponse
func (rs EmissionResource) ListByYear(w http.ResponseWriter, r *http.Request) {
	response := make(map[int]airPollutionEmissionsResponse)
	for year, emissions := range rs.Storage.FindAllByYears() {
		response[year] = newAirPollutionEmissionsResponse(emissions)
	}

	_ = json.NewEncoder(w).Encode(response)
}

// GetByYear TODO
// @Summary TODO
// @Description TODO
// @Tags emission year
// @Produce json
// @Router /emissions/year/{year} [get]
// @Param year path string true "year"
// @Success 200 {object} AirPollutionEmissionsResponse
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

	w.Header().Set("Content-Type", "application/json")
	response := newAirPollutionEmissionsResponse(yearEmissions)
	_ = render.Render(w, r, response)
}

// ListByCountry TODO
// @Summary TODO
// @Description TODO
// @Tags emission country
// @Produce json
// @Router /emissions/country/ [get]
// @Success 200 {object} map[year]AirPollutionEmissionsResponse
func (rs EmissionResource) ListByCountry(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]airPollutionEmissionsResponse)
	for country, emissions := range rs.Storage.FindAllByCountries() {
		response[country] = newAirPollutionEmissionsResponse(emissions)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// GetByCountry TODO
// @Summary TODO
// @Description TODO
// @Tags emission country
// @Produce json
// @Router /emissions/country/{name} [get]
// @Param name path string true "name of the country"
// @Success 200 {object} AirPollutionEmissionsResponse
// @Failure 400 {object} ErrResponse
// @Failure 404 {object} ErrResponse
func (rs EmissionResource) GetByCountry(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New(fmt.Sprintf("Country name missing"))))
		return
	}

	var countryEmissions []*model.Emissions
	for _, emissions := range rs.Storage.FindAllByCountry(name) {
		countryEmissions = append(countryEmissions, emissions)
	}

	response := newAirPollutionEmissionsResponse(countryEmissions)
	_ = render.Render(w, r, response)
}
