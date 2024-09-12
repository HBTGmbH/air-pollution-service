package model

type Country struct {
	Name string
	Code string
}

type Emissions struct {
	NOxEmissions   float64
	SO2Emissions   float64
	COEmissions    float64
	OCEmissions    float64
	NMVOCEmissions float64
	BCEmissions    float64
	NH3Emissions   float64
}
