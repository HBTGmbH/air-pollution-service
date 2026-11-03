// Package store provides in-memory access to the air pollution data set.
package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HBTGmbH/air-pollution-service/internal/csv"
	"github.com/HBTGmbH/air-pollution-service/internal/model"
)

// ErrDuplicateEmissions is returned by New when the data set contains more
// than one row for the same country and year.
var ErrDuplicateEmissions = errors.New("duplicate emissions")

// Store holds the complete data set in memory. Countries are identified by
// their lower-case ISO code or, if they have none, by their lower-case name
// with spaces replaced by dashes.
type Store struct {
	emissions map[string]map[int]model.Emissions
	countries map[string]model.Country
}

// New loads the complete data set from the given CSV file.
func New(file *csv.File) (*Store, error) {
	rows, err := file.ReadRows()
	if err != nil {
		return nil, fmt.Errorf("loading data set: %w", err)
	}

	emissions, err := toCountryEmissions(rows)
	if err != nil {
		return nil, err
	}

	return &Store{
		emissions: emissions,
		countries: toCountries(rows),
	}, nil
}

// FindAllByYears returns the emissions of all countries, grouped by year.
func (s *Store) FindAllByYears() map[int][]*model.Emissions {
	emissions := make(map[int][]*model.Emissions)
	for _, countryEmissions := range s.emissions {
		for year, countryEmissionsOfYear := range countryEmissions {
			emissions[year] = append(emissions[year], &countryEmissionsOfYear)
		}
	}

	return emissions
}

// FindAllByYear returns the emissions of all countries in the given year,
// keyed by country ID.
func (s *Store) FindAllByYear(year int) map[string]*model.Emissions {
	emissions := make(map[string]*model.Emissions)
	for id, countryEmissions := range s.emissions {
		countryEmissionsOfYear, found := countryEmissions[year]
		if found {
			emissions[id] = &countryEmissionsOfYear
		}
	}

	return emissions
}

// FindAllByCountries returns the emissions of all years, grouped by country ID.
func (s *Store) FindAllByCountries() map[string][]*model.Emissions {
	emissions := make(map[string][]*model.Emissions)
	for id, countryEmissions := range s.emissions {
		emissions[id] = make([]*model.Emissions, 0, len(countryEmissions))
		for _, countryEmissionsOfYear := range countryEmissions {
			emissions[id] = append(emissions[id], &countryEmissionsOfYear)
		}
	}

	return emissions
}

// FindAllByCountry returns the emissions of the given country, keyed by year.
func (s *Store) FindAllByCountry(id string) map[int]*model.Emissions {
	emissions := make(map[int]*model.Emissions)
	for year, countryEmissionsOfYear := range s.emissions[id] {
		emissions[year] = &countryEmissionsOfYear
	}

	return emissions
}

// GetCountry returns the country with the given ID, or nil if there is none.
func (s *Store) GetCountry(id string) *model.Country {
	country, exists := s.countries[id]
	if exists {
		return &country
	}

	return nil
}

// GetCountries returns all countries.
func (s *Store) GetCountries() []*model.Country {
	countries := make([]*model.Country, 0, len(s.countries))
	for _, country := range s.countries {
		countries = append(countries, &country)
	}

	return countries
}

func toCountryID(row *csv.Row) string {
	var id string
	if row.Code == "" {
		id = strings.ToLower(row.Entity)
	} else {
		id = strings.ToLower(row.Code)
	}
	id = strings.ReplaceAll(id, " ", "-")

	return id
}

func toCountryEmissions(rows []*csv.Row) (map[string]map[int]model.Emissions, error) {
	emissions := make(map[string]map[int]model.Emissions)
	for _, row := range rows {
		id := toCountryID(row)

		_, exists := emissions[id]
		if !exists {
			emissions[id] = make(map[int]model.Emissions)
		}

		_, exists = emissions[id][row.Year]
		if exists {
			return nil, fmt.Errorf("%w for year %d and country %s and code %s",
				ErrDuplicateEmissions, row.Year, row.Entity, row.Code)
		}

		emissions[id][row.Year] = model.Emissions{
			NOxEmissions:   row.NOxEmissions,
			SO2Emissions:   row.SO2Emissions,
			COEmissions:    row.COEmissions,
			OCEmissions:    row.OCEmissions,
			NMVOCEmissions: row.NMVOCEmissions,
			BCEmissions:    row.BCEmissions,
			NH3Emissions:   row.NH3Emissions,
		}
	}

	return emissions, nil
}

func toCountries(rows []*csv.Row) map[string]model.Country {
	countries := make(map[string]model.Country)
	for _, row := range rows {
		id := toCountryID(row)
		_, exists := countries[id]
		if exists {
			continue
		}
		countries[id] = model.Country{
			ID:   id,
			Name: row.Entity,
			Code: row.Code,
		}
	}

	return countries
}
