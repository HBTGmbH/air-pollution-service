package resource

import (
	"air-pollution-service/internal/model"
	"air-pollution-service/internal/repository"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type EmissionResource struct {
	*repository.Repository
}

type airPollutionResponse struct {
	Average  float64 `json:"average"`
	Median   float64 `json:"median"`
	Variance float64 `json:"variance"`
}

func newAirPollutionResponse(emissions []*model.Emissions, f func(emission *model.Emissions) float64) airPollutionResponse {
	total := 0.0

	for _, emission := range emissions {
		current := f(emission)
		total = total + current
	}

	return airPollutionResponse{
		Average:  total / float64(len(emissions)),
		Median:   0.0, // TODO
		Variance: 0.0, // TODO
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
	})

	return r
}

func (rs EmissionResource) ListByYear(w http.ResponseWriter, r *http.Request) {
	response := make(map[int]airPollutionEmissionsResponse)
	for year, emissions := range rs.FindAllByYears() {
		response[year] = newAirPollutionEmissionsResponse(emissions)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		render.Status(r, 500)
	}
}

func (rs EmissionResource) ListByCountry(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
	// TODO
}
