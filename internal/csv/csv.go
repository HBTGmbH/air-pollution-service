package csv

import (
	"github.com/gocarina/gocsv"
	"os"
)

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

type File struct {
	fileName string
}

func New(fileName string) *File {
	return &File{
		fileName: fileName,
	}
}

func (f *File) ReadRows() ([]*Row, error) {
	csvFile, err := os.OpenFile(f.fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			panic(err)
		}
	}(csvFile)

	var rows []*Row
	if err := gocsv.UnmarshalFile(csvFile, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}
