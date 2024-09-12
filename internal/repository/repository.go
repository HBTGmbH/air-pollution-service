package repository

import (
	"air-pollution-service/internal/csv"
	"air-pollution-service/internal/model"
	"fmt"
)

type Repository struct {
	emissions map[string]map[int]model.Emissions
	countries map[string]model.Country
}

func New(file *csv.File) (*Repository, error) {
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

	return &Repository{
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

func (r *Repository) FindAllByYears() map[int][]*model.Emissions {
	emissions := make(map[int][]*model.Emissions)
	for _, countryEmissions := range r.emissions {
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

func (r *Repository) FindAllByCountry(name string) map[int]*model.Emissions {
	emissions := make(map[int]*model.Emissions)
	for year, countryEmissionsOfYear := range r.emissions[name] {
		emissions[year] = &countryEmissionsOfYear
	}
	return emissions
}

func (r *Repository) GetCountry(name string) *model.Country {
	country, exists := r.countries[name]
	if exists {
		return &country
	}
	return nil
}

func (r *Repository) GetCountries() []*model.Country {
	var countries []*model.Country
	for _, country := range r.countries {
		countries = append(countries, &country)
	}
	return countries
}
