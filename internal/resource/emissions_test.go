package resource_test

import (
	"context"
	"encoding/json"
	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/HBTGmbH/air-pollution-service/internal/resource"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type fakeEmissionsStorage struct {
	emissions []*model.Emissions
}

func (s fakeEmissionsStorage) FindAllByYears() map[int][]*model.Emissions {
	result := make(map[int][]*model.Emissions)
	result[1234] = s.emissions
	return result
}

func (s fakeEmissionsStorage) FindAllByYear(year int) map[string]*model.Emissions {
	result := make(map[string]*model.Emissions)
	for i, e := range s.emissions {
		result[strconv.Itoa(i)] = e
	}
	return result
}

func (s fakeEmissionsStorage) FindAllByCountries() map[string][]*model.Emissions {
	result := make(map[string][]*model.Emissions)
	result["test"] = s.emissions
	return result
}

func (s fakeEmissionsStorage) FindAllByCountry(name string) map[int]*model.Emissions {
	result := make(map[int]*model.Emissions)
	for i, e := range s.emissions {
		result[i] = e
	}
	return result
}

func (s fakeEmissionsStorage) GetCountry(name string) *model.Country {
	panic("not implemented")
}

func (s fakeEmissionsStorage) GetCountries() []*model.Country {
	panic("not implemented")
}

func TestEmissionsListByYear(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	emissionsHandler := resource.EmissionResource{Storage: fakeEmissionsStorage{[]*model.Emissions{{
		NOxEmissions: 1,
	}, {
		NOxEmissions: 2,
	}, {
		NOxEmissions: 3,
	}}}}

	emissionsHandler.ListByYear(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// TODO validate response body
}

func TestEmissionsListByCountry(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	emissionsHandler := resource.EmissionResource{Storage: fakeEmissionsStorage{[]*model.Emissions{{
		NOxEmissions: 1,
	}, {
		NOxEmissions: 2,
	}, {
		NOxEmissions: 3,
	}}}}

	emissionsHandler.ListByCountry(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// TODO validate response body
}

func TestEmissionsGetByCountry(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	ctx.URLParams.Add("id", "Schlaraffenland")
	emissionsHandler := resource.EmissionResource{Storage: fakeEmissionsStorage{[]*model.Emissions{{
		NOxEmissions: 10,
	}, {
		NOxEmissions: 2,
	}, {
		NOxEmissions: 3,
	}}}}

	emissionsHandler.GetByCountry(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	airPollutionEmissions := resource.AirPollutionResponse{}
	err = json.Unmarshal(data, &airPollutionEmissions)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, airPollutionEmissions.NOxEmissions.Average)
	assert.Equal(t, 3.0, airPollutionEmissions.NOxEmissions.Median)
	assert.Equal(t, 3.56, airPollutionEmissions.NOxEmissions.StandardDeviation)
}

func TestEmissionsGetByYear(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	ctx.URLParams.Add("year", "666")
	emissionsHandler := resource.EmissionResource{Storage: fakeEmissionsStorage{[]*model.Emissions{{
		NOxEmissions: 10,
	}, {
		NOxEmissions: 2,
	}, {
		NOxEmissions: 3,
	}}}}

	emissionsHandler.GetByYear(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	airPollutionEmissions := resource.AirPollutionResponse{}
	err = json.Unmarshal(data, &airPollutionEmissions)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, airPollutionEmissions.NOxEmissions.Average)
	assert.Equal(t, 3.0, airPollutionEmissions.NOxEmissions.Median)
	assert.Equal(t, 3.56, airPollutionEmissions.NOxEmissions.StandardDeviation)
}
