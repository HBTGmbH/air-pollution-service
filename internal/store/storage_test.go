package store_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HBTGmbH/air-pollution-service/internal/csv"
	"github.com/HBTGmbH/air-pollution-service/internal/model"
	"github.com/HBTGmbH/air-pollution-service/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const header = "Entity,Code,Year,Nitrogen oxide (NOx),Sulphur dioxide (SO₂) emissions," +
	"Carbon monoxide (CO) emissions,Organic carbon (OC) emissions," +
	"Non-methane volatile organic compounds (NMVOC) emissions," +
	"Black carbon (BC) emissions,Ammonia (NH₃) emissions\n"

func writeFile(t *testing.T, rows ...string) *csv.File {
	t.Helper()

	fileName := filepath.Join(t.TempDir(), "air-pollution.csv")
	content := header + strings.Join(rows, "\n") + "\n"
	require.NoError(t, os.WriteFile(fileName, []byte(content), 0o600))

	return csv.New(fileName)
}

func TestStore(t *testing.T) {
	t.Parallel()

	s, err := store.New(writeFile(t,
		"Germany,DEU,2000,1,0,0,0,0,0,0",
		"Germany,DEU,2001,2,0,0,0,0,0,0",
		"North Africa,,2000,3,0,0,0,0,0,0",
	))
	require.NoError(t, err)

	germany := &model.Country{ID: "deu", Name: "Germany", Code: "DEU"}
	northAfrica := &model.Country{ID: "north-africa", Name: "North Africa"}
	germany2000 := &model.Emissions{NOxEmissions: 1}
	germany2001 := &model.Emissions{NOxEmissions: 2}
	northAfrica2000 := &model.Emissions{NOxEmissions: 3}

	t.Run("GetCountry", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, germany, s.GetCountry(germany.ID))
		assert.Equal(t, northAfrica, s.GetCountry(northAfrica.ID))
		assert.Nil(t, s.GetCountry("atlantis"))
	})

	t.Run("GetCountries", func(t *testing.T) {
		t.Parallel()
		assert.ElementsMatch(t, []*model.Country{germany, northAfrica}, s.GetCountries())
	})

	t.Run("FindAllByCountry", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t,
			map[int]*model.Emissions{2000: germany2000, 2001: germany2001},
			s.FindAllByCountry(germany.ID))
		assert.Empty(t, s.FindAllByCountry("atlantis"))
	})

	t.Run("FindAllByYear", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t,
			map[string]*model.Emissions{germany.ID: germany2000, northAfrica.ID: northAfrica2000},
			s.FindAllByYear(2000))
		assert.Empty(t, s.FindAllByYear(1999))
	})

	t.Run("FindAllByCountries", func(t *testing.T) {
		t.Parallel()
		byCountry := s.FindAllByCountries()
		assert.Len(t, byCountry, 2)
		assert.ElementsMatch(t, []*model.Emissions{germany2000, germany2001}, byCountry[germany.ID])
		assert.ElementsMatch(t, []*model.Emissions{northAfrica2000}, byCountry[northAfrica.ID])
	})

	t.Run("FindAllByYears", func(t *testing.T) {
		t.Parallel()
		byYear := s.FindAllByYears()
		assert.Len(t, byYear, 2)
		assert.ElementsMatch(t, []*model.Emissions{germany2000, northAfrica2000}, byYear[2000])
		assert.ElementsMatch(t, []*model.Emissions{germany2001}, byYear[2001])
	})
}

func TestStoreDuplicateEmissions(t *testing.T) {
	t.Parallel()

	_, err := store.New(writeFile(t,
		"Germany,DEU,2000,1,0,0,0,0,0,0",
		"Germany,DEU,2000,2,0,0,0,0,0,0",
	))

	require.ErrorIs(t, err, store.ErrDuplicateEmissions)
}

func TestStoreMissingFile(t *testing.T) {
	t.Parallel()

	_, err := store.New(csv.New(filepath.Join(t.TempDir(), "missing.csv")))

	require.ErrorIs(t, err, fs.ErrNotExist)
}
