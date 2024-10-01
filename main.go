package main

import (
	"air-pollution-service/config"
	_ "air-pollution-service/docs"
	"air-pollution-service/internal/csv"
	"air-pollution-service/internal/resource"
	"air-pollution-service/internal/store"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// get the config
	conf := config.New()

	// load raw data
	repo, err := store.New(csv.New(conf.AirPollutionFile))
	if err != nil {
		log.Panicf("Unable to load raw data from %s: %s", conf.AirPollutionFile, err)
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

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/swagger/*", httpSwagger.Handler())
	r.Mount("/countries", resource.CountryResource{Storage: repo}.Routes())
	r.Mount("/emissions", resource.EmissionResource{Storage: repo}.Routes())

	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: h2c.NewHandler(r, h2s),
	}

	// start the server
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("Failed to start server: %s", err)
		}
	}()

	log.Printf("Server started sucessfully!")
	sig := <-c
	if sig == syscall.SIGTERM {
		// Build "sleep before shutdown" into the application to handle deployment
		// scenarios where the service might still receive traffic after a shutdown
		// such as in a Kubernetes rolling deployment/update. This will allow the
		// service to gracefully shutdown after a configurable duration after we can
		// be sure that all traffic has been shifted to the new version.
		// See also: https://learnk8s.io/graceful-shutdown
		parsedSleepDuration, err := time.ParseDuration(conf.Server.SleepDurationBeforeShutdown)
		if err != nil {
			log.Printf("Failed to parse SLEEP_DURATION_BEFORE_SHUTDOWN duration %s. Won't sleep: %s", conf.Server.SleepDurationBeforeShutdown, err)
		} else {
			log.Printf("Waiting for %s before shutting down...", parsedSleepDuration.String())
			time.Sleep(parsedSleepDuration)
		}
	}
	log.Printf("Shutting down server gracefully...")
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Panicf("Failed to shutdown server: %s", err)
	}
	log.Printf("Server shutdown successfully. Exiting.")
}
