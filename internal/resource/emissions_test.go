package resource_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HBTGmbH/air-pollution-service/internal/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmissionsListByYear(t *testing.T) {
	t.Parallel()

	handler := resource.EmissionResource{Storage: fakeStorage{emissions: noxEmissions(1, 2, 3)}}

	w := serve(t, handler.ListByYear, nil)

	require.Equal(t, http.StatusOK, w.Code)
	var response resource.AirPollutionYears
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.Contains(t, response.Years, 1234)
	assert.Equal(t,
		resource.EmissionsResponse{Average: 2, Median: 2, StandardDeviation: 0.82},
		response.Years[1234].NOxEmissions)
}

func TestEmissionsListByCountry(t *testing.T) {
	t.Parallel()

	handler := resource.EmissionResource{Storage: fakeStorage{emissions: noxEmissions(1, 2, 3)}}

	w := serve(t, handler.ListByCountry, nil)

	require.Equal(t, http.StatusOK, w.Code)
	var response resource.AirPollutionCountries
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.Contains(t, response.Countries, "test")
	assert.Equal(t,
		resource.EmissionsResponse{Average: 2, Median: 2, StandardDeviation: 0.82},
		response.Countries["test"].NOxEmissions)
}

func TestEmissionsGet(t *testing.T) {
	t.Parallel()

	endpoints := map[string]struct {
		handler func(resource.EmissionResource, http.ResponseWriter, *http.Request)
		params  map[string]string
	}{
		"by year":    {resource.EmissionResource.GetByYear, map[string]string{"year": "2000"}},
		"by country": {resource.EmissionResource.GetByCountry, map[string]string{"id": "DEU"}},
	}
	for name, endpoint := range endpoints {
		serveEndpoint := func(t *testing.T, storage fakeStorage) *httptest.ResponseRecorder {
			t.Helper()
			handler := resource.EmissionResource{Storage: storage}

			return serve(t, func(w http.ResponseWriter, r *http.Request) {
				endpoint.handler(handler, w, r)
			}, endpoint.params)
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := serveEndpoint(t, fakeStorage{emissions: noxEmissions(10, 2, 3)})

			require.Equal(t, http.StatusOK, w.Code)
			var response resource.AirPollutionResponse
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
			assert.Equal(t,
				resource.EmissionsResponse{Average: 5, Median: 3, StandardDeviation: 3.56},
				response.NOxEmissions)
		})

		// Without any emissions there are no statistics to compute.
		t.Run(name+" without data", func(t *testing.T) {
			t.Parallel()

			w := serveEndpoint(t, fakeStorage{})

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestEmissionsGetByYearInvalid(t *testing.T) {
	t.Parallel()

	handler := resource.EmissionResource{Storage: fakeStorage{emissions: noxEmissions(1)}}

	w := serve(t, handler.GetByYear, map[string]string{"year": "last year"})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
