package resource

import (
	"air-pollution-service/internal/model"
	"air-pollution-service/internal/store"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/montanaflynn/stats"
	"net/http"
)

type EmissionResource struct {
	Storage store.Storage
}

type airPollutionResponse struct {
	Average           float64 `json:"average"`
	Median            float64 `json:"median"`
	StandardDeviation float64 `json:"standard_deviation"`
}

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

type airPollutionEmissionsResponse struct {
	NOxEmissions   airPollutionResponse `json:"nox_emissions"`
	SO2Emissions   airPollutionResponse `json:"sulphur_dioxide_emissions"`
	COEmissions    airPollutionResponse `json:"carbon_monoxide_emissions"`
	OCEmissions    airPollutionResponse `json:"organic_carbon_emissions"`
	NMVOCEmissions airPollutionResponse `json:"nmvoc_emissions"`
	BCEmissions    airPollutionResponse `json:"black_carbon_emissions"`
	NH3Emissions   airPollutionResponse `json:"ammonia_emissions"`
}

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
	})

	r.Route("/country/", func(r chi.Router) {
		r.Get("/", rs.ListByCountry)
		r.Get("/{name}", rs.GetByCountry)
	})

	return r
}

func (rs EmissionResource) ListByYear(w http.ResponseWriter, r *http.Request) {
	response := make(map[int]airPollutionEmissionsResponse)
	for year, emissions := range rs.Storage.FindAllByYears() {
		response[year] = newAirPollutionEmissionsResponse(emissions)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		render.Status(r, 500)
	}
}

func (rs EmissionResource) ListByCountry(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]airPollutionEmissionsResponse)
	for country, emissions := range rs.Storage.FindAllByCountries() {
		response[country] = newAirPollutionEmissionsResponse(emissions)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		render.Status(r, 500)
	}
}

func (rs EmissionResource) GetByCountry(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		if err := render.Render(w, r, ErrRender(fmt.Sprintf("Country name missing"), 400)); err != nil {
			render.Status(r, 500)
		}
		return
	}

	var countryEmissions []*model.Emissions
	for _, emissions := range rs.Storage.FindAllByCountry(name) {
		countryEmissions = append(countryEmissions, emissions)
	}

	response := newAirPollutionEmissionsResponse(countryEmissions)
	if err := render.Render(w, r, response); err != nil {
		render.Status(r, 500)
	}
}
