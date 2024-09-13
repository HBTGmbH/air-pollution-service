package store

import (
	"air-pollution-service/internal/csv"
	"air-pollution-service/internal/model"
	"fmt"
)

type Store struct {
	emissions map[string]map[int]model.Emissions
	countries map[string]model.Country
}

type Storage interface {
	FindAllByYears() map[int][]*model.Emissions
	FindAllByYear(year int) map[string]*model.Emissions
	FindAllByCountries() map[string][]*model.Emissions
	FindAllByCountry(name string) map[int]*model.Emissions
	GetCountry(name string) *model.Country
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

func toCountryEmissions(rowsFromFile []*csv.Row) (map[string]map[int]model.Emissions, error) {
	emissions := make(map[string]map[int]model.Emissions)
	for _, row := range rowsFromFile {
		_, exists := emissions[row.Entity]
		if !exists {
			emissions[row.Entity] = make(map[int]model.Emissions)
		}

		_, exists = emissions[row.Entity][row.Year]
		if exists {
			return nil, fmt.Errorf("duplicate emissions for year %d and country %s", row.Year, row.Entity)
		}

		emissions[row.Entity][row.Year] = model.Emissions{
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

func toCountries(rowsFromFile []*csv.Row) (map[string]model.Country, error) {
	countries := make(map[string]model.Country)
	for _, row := range rowsFromFile {
		_, exists := countries[row.Entity]
		if exists {
			continue
		}
		countries[row.Entity] = model.Country{
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
	for name, countryEmissions := range s.emissions {
		countryEmissionsOfYear, found := countryEmissions[year]
		if found {
			emissions[name] = &countryEmissionsOfYear
		}
	}
	return emissions
}

func (s *Store) FindAllByCountries() map[string][]*model.Emissions {
	emissions := make(map[string][]*model.Emissions)
	for name, countryEmissions := range s.emissions {
		emissions[name] = []*model.Emissions{}
		for _, countryEmissionsOfYear := range countryEmissions {
			emissions[name] = append(emissions[name], &countryEmissionsOfYear)
		}
	}
	return emissions
}

func (s *Store) FindAllByCountry(name string) map[int]*model.Emissions {
	emissions := make(map[int]*model.Emissions)
	for year, countryEmissionsOfYear := range s.emissions[name] {
		emissions[year] = &countryEmissionsOfYear
	}
	return emissions
}

func (s *Store) GetCountry(name string) *model.Country {
	country, exists := s.countries[name]
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
