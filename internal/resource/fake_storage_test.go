package resource_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/go-chi/chi/v5"
)

// fakeStorage serves the same emissions for every year and country.
type fakeStorage struct {
	countries []*model.Country
	emissions []*model.Emissions
}

func (s fakeStorage) FindAllByYears() map[int][]*model.Emissions {
	return map[int][]*model.Emissions{1234: s.emissions}
}

func (s fakeStorage) FindAllByYear(_ int) map[string]*model.Emissions {
	result := make(map[string]*model.Emissions, len(s.emissions))
	for i, e := range s.emissions {
		result[strconv.Itoa(i)] = e
	}

	return result
}

func (s fakeStorage) FindAllByCountries() map[string][]*model.Emissions {
	return map[string][]*model.Emissions{"test": s.emissions}
}

func (s fakeStorage) FindAllByCountry(_ string) map[int]*model.Emissions {
	result := make(map[int]*model.Emissions, len(s.emissions))
	for i, e := range s.emissions {
		result[i] = e
	}

	return result
}

func (s fakeStorage) GetCountry(_ string) *model.Country {
	if len(s.countries) == 0 {
		return nil
	}

	return s.countries[0]
}

func (s fakeStorage) GetCountries() []*model.Country {
	return s.countries
}

// noxEmissions returns one emissions entry per value, with only the NOx
// emissions set.
func noxEmissions(values ...float64) []*model.Emissions {
	emissions := make([]*model.Emissions, 0, len(values))
	for _, value := range values {
		emissions = append(emissions, &model.Emissions{NOxEmissions: value})
	}

	return emissions
}

// serve calls handler with a GET request carrying the given chi URL parameters.
func serve(t *testing.T, handler http.HandlerFunc, params map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	routeContext := chi.NewRouteContext()
	for key, value := range params {
		routeContext.URLParams.Add(key, value)
	}
	ctx := context.WithValue(t.Context(), chi.RouteCtxKey, routeContext)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	return w
}
