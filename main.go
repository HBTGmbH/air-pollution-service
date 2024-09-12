package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gocarina/gocsv"
	"log"
	"net/http"
	"os"
)

type Emissions struct {
	Entity         string  `csv:"Entity"`
	Code           string  `csv:"Code"`
	Year           int     `csv:"Year"`
	NOxEmissions   float32 `csv:"Nitrogen oxide (NOx)"`
	SO2Emissions   float32 `csv:"Sulphur dioxide (SO₂) emissions"`
	COEmissions    float32 `csv:"Carbon monoxide (CO) emissions"`
	OCEmissions    float32 `csv:"Organic carbon (OC) emissions"`
	NMVOCEmissions float32 `csv:"Non-methane volatile organic compounds (NMVOC) emissions"`
	BCEmissions    float32 `csv:"Black carbon (BC) emissions"`
	NH3Emissions   float32 `csv:"Ammonia (NH₃) emissions"`
}

func ReadCsv() []*Emissions {
	// Try to open the example.csv file in read-write mode.
	csvFile, csvFileError := os.OpenFile("air-pollution.csv", os.O_RDWR, os.ModePerm)
	// If an error occurs during os.OpenFIle, panic and halt execution.
	if csvFileError != nil {
		panic(csvFileError)
	}
	// Ensure the file is closed once the function returns
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			panic(err)
		}
	}(csvFile)

	var emissions []*Emissions
	// Parse the CSV data into the emissions slice. If an error occurs, panic.
	if unmarshalError := gocsv.UnmarshalFile(csvFile, &emissions); unmarshalError != nil {
		panic(unmarshalError)
	}

	return emissions
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/countries", GetCountries)
	r.Get("/countries/{id}", GetCountry)
	r.Get("/countries/{id}/{year}", GetCountryByYear)
	r.Get("/years", GetYears)
	r.Get("/years/{year}", GetYear)
	r.Get("/years/{year}/{id}", GetYearByCountry)

	log.Println("Starting server on :3000")
	http.ListenAndServe(":3000", r)
}

func GetCountries(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func GetCountry(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func GetCountryByYear(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func GetYears(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func GetYear(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}

func GetYearByCountry(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
}
