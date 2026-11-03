package resource_test

import (
	"net/http"
	"testing"

	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/HBTGmbH/air-pollution-service/internal/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var schlaraffenland = &model.Country{ID: "sch", Name: "Schlaraffenland", Code: "SCH"}

const schlaraffenlandJSON = `{"id": "sch", "name": "Schlaraffenland", "code": "SCH"}`

func TestCountriesList(t *testing.T) {
	t.Parallel()

	handler := resource.CountryResource{Storage: fakeStorage{countries: []*model.Country{
		schlaraffenland,
		{ID: "utopia", Name: "Utopia"},
	}}}

	w := serve(t, handler.List, nil)

	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t,
		`[`+schlaraffenlandJSON+`, {"id": "utopia", "name": "Utopia", "code": ""}]`,
		w.Body.String())
}

func TestCountriesListEmpty(t *testing.T) {
	t.Parallel()

	handler := resource.CountryResource{Storage: fakeStorage{}}

	w := serve(t, handler.List, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountriesGetNonExisting(t *testing.T) {
	t.Parallel()

	handler := resource.CountryResource{Storage: fakeStorage{}}

	w := serve(t, handler.Get, map[string]string{"id": "atlantis"})

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCountriesGetExisting(t *testing.T) {
	t.Parallel()

	handler := resource.CountryResource{Storage: fakeStorage{countries: []*model.Country{schlaraffenland}}}

	w := serve(t, handler.Get, map[string]string{"id": schlaraffenland.Name})

	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, schlaraffenlandJSON, w.Body.String())
}
