// Package model contains the domain types of the air pollution data set.
package model

// Country is a country or region of the data set.
type Country struct {
	ID   string
	Name string
	Code string
}

// Emissions are the emissions of a single country in a single year.
type Emissions struct {
	NOxEmissions   float64
	SO2Emissions   float64
	COEmissions    float64
	OCEmissions    float64
	NMVOCEmissions float64
	BCEmissions    float64
	NH3Emissions   float64
}
