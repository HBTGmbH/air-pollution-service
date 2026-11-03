// Command air-pollution-service serves statistics about historical air
// pollution data of all countries via a REST API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/HBTGmbH/air-pollution-service/config"
	_ "github.com/HBTGmbH/air-pollution-service/docs"
	"github.com/HBTGmbH/air-pollution-service/internal/csv"
	"github.com/HBTGmbH/air-pollution-service/internal/drain"
	"github.com/HBTGmbH/air-pollution-service/internal/resource"
	"github.com/HBTGmbH/air-pollution-service/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpswagger "github.com/swaggo/http-swagger"
)

var (
	draining             atomic.Bool
	openConnectionsCount atomic.Int64
)

func main() {
	// get the config
	conf, err := config.New()
	if err != nil {
		log.Fatalf("Invalid configuration: %s", err)
	}

	// load raw data
	repo, err := store.New(csv.New(conf.AirPollutionFile))
	if err != nil {
		log.Panicf("Unable to load raw data from %s: %s", conf.AirPollutionFile, err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// build API router
	r := chi.NewRouter()
	r.Use(drain.WithConnectionDraining(draining.Load))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/swagger/*", httpswagger.Handler())
	r.Mount("/countries", resource.CountryResource{Storage: repo}.Routes())
	r.Mount("/emissions", resource.EmissionResource{Storage: repo}.Routes())

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Server.Port),
		Handler:           r,
		IdleTimeout:       conf.Server.IdleTimeout,
		ReadHeaderTimeout: conf.Server.ReadHeaderTimeout,
		ConnState: func(_ net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				openConnectionsCount.Add(1)
			case http.StateClosed, http.StateHijacked:
				openConnectionsCount.Add(-1)
			default:
				// no-op
			}
		},
	}

	// start the server
	go func() {
		serveErr := server.ListenAndServe()
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Panicf("Failed to start server: %s", serveErr)
		}
	}()

	log.Print("Server started successfully!")
	<-c
	gracefulShutdown(server, conf)
}

func gracefulShutdown(server *http.Server, conf *config.Conf) {
	draining.Store(true)
	if conf.Server.ShutdownSleepDuration > 0 {
		log.Printf("Waiting for %s until no new connections should come in anymore.",
			conf.Server.ShutdownSleepDuration.String())
		time.Sleep(conf.Server.ShutdownSleepDuration)
	}
	if conf.Server.DrainTimeout > 0 {
		log.Printf("Draining idle connections for up to %s.", conf.Server.DrainTimeout.String())
		drainStart := time.Now()
		conns := openConnectionsCount.Load()
		for conns > 0 && time.Since(drainStart) < conf.Server.DrainTimeout {
			log.Printf("Waiting for connections to drain. Remaining connections: %d", conns)
			time.Sleep(1 * time.Second)
			conns = openConnectionsCount.Load()
		}
		log.Printf("Finished draining connections. Remaining connections: %d", conns)
	}
	log.Printf("Shutting down server with a timeout of %s.", conf.Server.ShutdownTimeout.String())
	ctx, cancel := context.WithTimeout(context.Background(), conf.Server.ShutdownTimeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error during server shutdown: %s", err)
	} else {
		log.Print("Server shutdown completed successfully.")
	}
}
