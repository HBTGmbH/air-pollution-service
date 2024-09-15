package resource

import (
	"air-pollution-service/internal/model"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeCountryStorage struct {
	countries []*model.Country
}

func (s fakeCountryStorage) FindAllByYears() map[int][]*model.Emissions {
	panic("not implemented")
}

func (s fakeCountryStorage) FindAllByYear(year int) map[string]*model.Emissions {
	panic("not implemented")
}

func (s fakeCountryStorage) FindAllByCountries() map[string][]*model.Emissions {
	panic("not implemented")
}

func (s fakeCountryStorage) FindAllByCountry(name string) map[int]*model.Emissions {
	panic("not implemented")
}

func (s fakeCountryStorage) GetCountry(name string) *model.Country {
	if len(s.countries) == 0 {
		return nil
	}
	return s.countries[0]
}

func (s fakeCountryStorage) GetCountries() []*model.Country {
	return s.countries
}

func TestCountriesGetNonExisting(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	ctx.URLParams.Add("id", "Schlaraffenland")
	countryHandler := CountryResource{Storage: fakeCountryStorage{[]*model.Country{}}}

	countryHandler.Get(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func TestCountriesGetExisting(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	ctx.URLParams.Add("id", "Schlaraffenland")
	countryHandler := CountryResource{Storage: fakeCountryStorage{[]*model.Country{{
		Name: "Schlaraffenland",
		Code: "SCH",
		Id:   "sch",
	}}}}

	countryHandler.Get(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	country := countryResponse{}
	err = json.Unmarshal(data, &country)
	assert.Nil(t, err)
	assert.Equal(t, "Schlaraffenland", country.Name)
	assert.Equal(t, "SCH", country.Code)
	assert.Equal(t, "sch", country.Id)
}
