package main

import (
	"air-pollution-service/config"
	"air-pollution-service/internal/csv"
	"air-pollution-service/internal/repository"
	"air-pollution-service/internal/resource"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func main() {
	c := config.New()

	repo, err := repository.New(csv.New(c.AirPollutionFile))
	if err != nil {
		log.Panic(err)
	}

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/countries", resource.CountryResource{Repository: repo}.Routes())
	r.Mount("/emissions", resource.EmissionResource{Repository: repo}.Routes())

	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Server.Port), r)
	if err != nil {
		log.Panicf("Failed to start server %s", err)
	}
}
