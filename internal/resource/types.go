package resource

import (
	"air-pollution-service/internal/model"
	"github.com/go-chi/render"
	"github.com/montanaflynn/stats"
	"net/http"
)

// CountryResponse renderer type for country data
type countryResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
} // @name CountryResponse

func newCountryResponse(country *model.Country) countryResponse {
	return countryResponse{
		Id:   country.Id,
		Name: country.Name,
		Code: country.Code,
	}
}

func (hr countryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// EmissionsResponse renderer type for accumulated air pollution datasets
type EmissionsResponse struct {
	Average           float64 `json:"average"`
	Median            float64 `json:"median"`
	StandardDeviation float64 `json:"standard_deviation"`
} // @name EmissionsResponse

func newEmissionsResponse(emissions []*model.Emissions, f func(emission *model.Emissions) float64) EmissionsResponse {
	total := 0.0
	values := make([]float64, len(emissions))

	for i, emission := range emissions {
		current := f(emission)
		values[i] = current
		total = total + current
	}

	median, _ := stats.Median(values)
	standardDeviation, _ := stats.StandardDeviation(values)

	return EmissionsResponse{
		Average:           total / float64(len(emissions)),
		Median:            median,
		StandardDeviation: standardDeviation,
	}
}

// AirPollutionResponse renderer type for a single air pollution dataset
type AirPollutionResponse struct {
	NOxEmissions   EmissionsResponse `json:"nox_emissions"`
	SO2Emissions   EmissionsResponse `json:"sulphur_dioxide_emissions"`
	COEmissions    EmissionsResponse `json:"carbon_monoxide_emissions"`
	OCEmissions    EmissionsResponse `json:"organic_carbon_emissions"`
	NMVOCEmissions EmissionsResponse `json:"nmvoc_emissions"`
	BCEmissions    EmissionsResponse `json:"black_carbon_emissions"`
	NH3Emissions   EmissionsResponse `json:"ammonia_emissions"`
} // @name AirPollutionResponse

func newAirPollutionEmissionsResponse(emissions []*model.Emissions) AirPollutionResponse {
	return AirPollutionResponse{
		NOxEmissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NOxEmissions
		}),
		SO2Emissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.SO2Emissions
		}),
		COEmissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.COEmissions
		}),
		OCEmissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.OCEmissions
		}),
		NMVOCEmissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NMVOCEmissions
		}),
		BCEmissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.BCEmissions
		}),
		NH3Emissions: newEmissionsResponse(emissions, func(emissions *model.Emissions) float64 {
			return emissions.NH3Emissions
		}),
	}
}

func (hr AirPollutionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// AirPollutionYears renderer type for air pollution data by year
type AirPollutionYears struct {
	Years map[int]AirPollutionResponse `json:"years"`
} // @name AirPollutionYears

func (hr AirPollutionYears) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// AirPollutionCountries renderer type for air pollution data by country
type AirPollutionCountries struct {
	Countries map[string]AirPollutionResponse `json:"countries"`
} // @name AirPollutionCountries

func (hr AirPollutionCountries) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	StatusCode int    `json:"-"`
	StatusText string `json:"text,omitempty"`
	Err        error  `json:"-"`
	ErrCode    int    `json:"code,omitempty" example:"404"`
	ErrText    string `json:"error,omitempty" example:"The requested resource was not found on the server."`
} // @name ErrResponse

// ErrRender returns a structured http response in case of rendering errors
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: http.StatusUnprocessableEntity,
		StatusText: "Error rendering response",
	}
}

// ErrInvalidRequest returns a structured http response in case of an invalid request
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request",
	}
}

// ErrNotFound returns a structured http response in case a resource was not found
func ErrNotFound() render.Renderer {
	return &ErrResponse{
		StatusCode: http.StatusNotFound,
		StatusText: "Resource not found",
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
