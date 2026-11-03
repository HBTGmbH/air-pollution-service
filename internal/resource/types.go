package resource

import (
	"fmt"
	"net/http"

	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/montanaflynn/stats"
)

// decimalPlaces is the precision of all statistics in API responses.
const decimalPlaces = 2

// CountryResponse renderer type for country data.
type CountryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
} //	@name	CountryResponse

func newCountryResponse(country *model.Country) CountryResponse {
	return CountryResponse{
		ID:   country.ID,
		Name: country.Name,
		Code: country.Code,
	}
}

// Render implements render.Renderer.
func (CountryResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// EmissionsResponse renderer type for accumulated air pollution datasets.
type EmissionsResponse struct {
	Average           float64 `json:"average"`
	Median            float64 `json:"median"`
	StandardDeviation float64 `json:"standard_deviation"`
} //	@name	EmissionsResponse

// newEmissionsResponse computes the statistics of one kind of emissions,
// selected by value, over all given emissions. It fails if emissions is empty.
func newEmissionsResponse(
	emissions []*model.Emissions,
	value func(*model.Emissions) float64,
) (EmissionsResponse, error) {
	values := make(stats.Float64Data, len(emissions))
	for i, emission := range emissions {
		values[i] = value(emission)
	}

	average, err := roundedStatistic(values, stats.Mean)
	if err != nil {
		return EmissionsResponse{}, fmt.Errorf("computing average: %w", err)
	}

	median, err := roundedStatistic(values, stats.Median)
	if err != nil {
		return EmissionsResponse{}, fmt.Errorf("computing median: %w", err)
	}

	standardDeviation, err := roundedStatistic(values, stats.StandardDeviation)
	if err != nil {
		return EmissionsResponse{}, fmt.Errorf("computing standard deviation: %w", err)
	}

	return EmissionsResponse{
		Average:           average,
		Median:            median,
		StandardDeviation: standardDeviation,
	}, nil
}

func roundedStatistic(
	values stats.Float64Data,
	statistic func(stats.Float64Data) (float64, error),
) (float64, error) {
	result, err := statistic(values)
	if err != nil {
		return 0, err
	}

	rounded, err := stats.Round(result, decimalPlaces)
	if err != nil {
		return 0, fmt.Errorf("rounding %v: %w", result, err)
	}

	return rounded, nil
}

// AirPollutionResponse renderer type for a single air pollution dataset.
type AirPollutionResponse struct {
	NOxEmissions   EmissionsResponse `json:"nox_emissions"`
	SO2Emissions   EmissionsResponse `json:"sulphur_dioxide_emissions"`
	COEmissions    EmissionsResponse `json:"carbon_monoxide_emissions"`
	OCEmissions    EmissionsResponse `json:"organic_carbon_emissions"`
	NMVOCEmissions EmissionsResponse `json:"nmvoc_emissions"`
	BCEmissions    EmissionsResponse `json:"black_carbon_emissions"`
	NH3Emissions   EmissionsResponse `json:"ammonia_emissions"`
} //	@name	AirPollutionResponse

// newAirPollutionResponse computes the statistics of all kinds of emissions
// over all given emissions. It fails if emissions is empty.
func newAirPollutionResponse(emissions []*model.Emissions) (AirPollutionResponse, error) {
	var response AirPollutionResponse
	fields := []struct {
		target *EmissionsResponse
		value  func(*model.Emissions) float64
	}{
		{&response.NOxEmissions, func(e *model.Emissions) float64 { return e.NOxEmissions }},
		{&response.SO2Emissions, func(e *model.Emissions) float64 { return e.SO2Emissions }},
		{&response.COEmissions, func(e *model.Emissions) float64 { return e.COEmissions }},
		{&response.OCEmissions, func(e *model.Emissions) float64 { return e.OCEmissions }},
		{&response.NMVOCEmissions, func(e *model.Emissions) float64 { return e.NMVOCEmissions }},
		{&response.BCEmissions, func(e *model.Emissions) float64 { return e.BCEmissions }},
		{&response.NH3Emissions, func(e *model.Emissions) float64 { return e.NH3Emissions }},
	}
	for _, field := range fields {
		emissionsResponse, err := newEmissionsResponse(emissions, field.value)
		if err != nil {
			return AirPollutionResponse{}, err
		}
		*field.target = emissionsResponse
	}

	return response, nil
}

// Render implements render.Renderer.
func (*AirPollutionResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// AirPollutionYears renderer type for air pollution data by year.
type AirPollutionYears struct {
	Years map[int]AirPollutionResponse `json:"years"`
} //	@name	AirPollutionYears

// Render implements render.Renderer.
func (AirPollutionYears) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// AirPollutionCountries renderer type for air pollution data by country.
type AirPollutionCountries struct {
	Countries map[string]AirPollutionResponse `json:"countries"`
} //	@name	AirPollutionCountries

// Render implements render.Renderer.
func (AirPollutionCountries) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
