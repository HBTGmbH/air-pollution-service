package store

import (
	"fmt"
	"strings"

	"github.com/HBTGmbH/air-pollution-service/internal/csv"
	"github.com/HBTGmbH/air-pollution-service/internal/model"
)

type Store struct {
	emissions map[string]map[int]model.Emissions
	countries map[string]model.Country
}

type Storage interface {
	FindAllByYears() map[int][]*model.Emissions
	FindAllByYear(year int) map[string]*model.Emissions
	FindAllByCountries() map[string][]*model.Emissions
	FindAllByCountry(id string) map[int]*model.Emissions
	GetCountry(id string) *model.Country
	GetCountries() []*model.Country
}

func New(file *csv.File) (*Store, error) {
	rows, err := file.ReadRows()
	if err != nil {
		return nil, err
	}

	emissions, err := toCountryEmissions(rows)
	if err != nil {
		return nil, err
	}

	countries, err := toCountries(rows)
	if err != nil {
		return nil, err
	}

	return &Store{
		emissions: emissions,
		countries: countries,
	}, nil
}

func toCountryId(row *csv.Row) string {
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
		id := toCountryId(row)

		_, exists := emissions[id]
		if !exists {
			emissions[id] = make(map[int]model.Emissions)
		}

		_, exists = emissions[id][row.Year]
		if exists {
			return nil, fmt.Errorf("duplicate emissions for year %d and country %s and code %s", row.Year, row.Entity, row.Code)
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

func toCountries(rows []*csv.Row) (map[string]model.Country, error) {
	countries := make(map[string]model.Country)
	for _, row := range rows {
		id := toCountryId(row)
		_, exists := countries[id]
		if exists {
			continue
		}
		countries[id] = model.Country{
			Id:   id,
			Name: row.Entity,
			Code: row.Code,
		}
	}
	return countries, nil
}

func (s *Store) FindAllByYears() map[int][]*model.Emissions {
	emissions := make(map[int][]*model.Emissions)
	for _, countryEmissions := range s.emissions {
		for year, countryEmissionsOfYear := range countryEmissions {
			_, found := emissions[year]
			if !found {
				emissions[year] = []*model.Emissions{&countryEmissionsOfYear}
			} else {
				emissions[year] = append(emissions[year], &countryEmissionsOfYear)
			}
		}
	}
	return emissions
}

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

func (s *Store) FindAllByCountries() map[string][]*model.Emissions {
	emissions := make(map[string][]*model.Emissions)
	for id, countryEmissions := range s.emissions {
		emissions[id] = []*model.Emissions{}
		for _, countryEmissionsOfYear := range countryEmissions {
			emissions[id] = append(emissions[id], &countryEmissionsOfYear)
		}
	}
	return emissions
}

func (s *Store) FindAllByCountry(id string) map[int]*model.Emissions {
	emissions := make(map[int]*model.Emissions)
	for year, countryEmissionsOfYear := range s.emissions[id] {
		emissions[year] = &countryEmissionsOfYear
	}
	return emissions
}

func (s *Store) GetCountry(id string) *model.Country {
	country, exists := s.countries[id]
	if exists {
		return &country
	}
	return nil
}

func (s *Store) GetCountries() []*model.Country {
	var countries []*model.Country
	for _, country := range s.countries {
		countries = append(countries, &country)
	}
	return countries
}
