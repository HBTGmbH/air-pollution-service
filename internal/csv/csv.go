// Package csv reads the raw air pollution data set from a CSV file.
package csv

import (
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

// Row is a single row of the air pollution data set: the emissions of one
// country in one year.
type Row struct {
	Entity         string  `csv:"Entity"`
	Code           string  `csv:"Code"`
	Year           int     `csv:"Year"`
	NOxEmissions   float64 `csv:"Nitrogen oxide (NOx)"`
	SO2Emissions   float64 `csv:"Sulphur dioxide (SO₂) emissions"`
	COEmissions    float64 `csv:"Carbon monoxide (CO) emissions"`
	OCEmissions    float64 `csv:"Organic carbon (OC) emissions"`
	NMVOCEmissions float64 `csv:"Non-methane volatile organic compounds (NMVOC) emissions"`
	BCEmissions    float64 `csv:"Black carbon (BC) emissions"`
	NH3Emissions   float64 `csv:"Ammonia (NH₃) emissions"`
}

// File is a CSV file containing the air pollution data set.
type File struct {
	fileName string
}

// New returns the CSV file with the given name.
func New(fileName string) *File {
	return &File{
		fileName: fileName,
	}
}

// ReadRows reads and parses all rows of the file.
func (f *File) ReadRows() ([]*Row, error) {
	log.Printf("Load data from %s ...", f.fileName)
	data, err := os.ReadFile(f.fileName)
	if err != nil {
		return nil, fmt.Errorf("reading CSV file: %w", err)
	}

	var rows []*Row
	err = gocsv.UnmarshalBytes(data, &rows)
	if err != nil {
		return nil, fmt.Errorf("parsing CSV file %s: %w", f.fileName, err)
	}

	return rows, nil
}
