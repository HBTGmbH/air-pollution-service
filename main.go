package main

import (
	"air-pollution-service/config"
	_ "air-pollution-service/docs"
	"air-pollution-service/internal/csv"
	"air-pollution-service/internal/resource"
	"air-pollution-service/internal/store"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// get the config
	conf := config.New()

	// load raw data
	repo, err := store.New(csv.New(conf.AirPollutionFile))
	if err != nil {
		log.Panicf("Unable to load raw data from %s", conf.AirPollutionFile)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// build API router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/swagger/*", httpSwagger.Handler())
	r.Mount("/countries", resource.CountryResource{Storage: repo}.Routes())
	r.Mount("/emissions", resource.EmissionResource{Storage: repo}.Routes())

	// start the server
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Server.Port), r); err != nil {
			log.Panicf("Failed to start server %s", err)
		}
	}()

	log.Printf("Server started sucessfully!")
	<-c
	log.Printf("Shutting down server!")
	os.Exit(0)
}
