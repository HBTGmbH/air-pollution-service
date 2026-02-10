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
	"github.com/HBTGmbH/air-pollution-service/internal/resource"
	"github.com/HBTGmbH/air-pollution-service/internal/store"
	"github.com/HBTGmbH/air-pollution-service/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	draining             = atomic.Bool{}
	openConnectionsCount int64
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
	r.Use(util.WithConnectionDraining(func() bool { return draining.Load() }))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r.Get("/swagger/*", httpSwagger.Handler())
	r.Mount("/countries", resource.CountryResource{Storage: repo}.Routes())
	r.Mount("/emissions", resource.EmissionResource{Storage: repo}.Routes())

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", conf.Server.Port),
		Handler:     r,
		IdleTimeout: conf.Server.IdleTimeout,
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				atomic.AddInt64(&openConnectionsCount, 1)
			case http.StateClosed, http.StateHijacked:
				atomic.AddInt64(&openConnectionsCount, -1)
			default:
				// no-op
			}
		},
	}

	// start the server
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("Failed to start server: %s", err)
		}
	}()

	log.Printf("Server started sucessfully!")
	_ = <-c
	gracefulShutdown(server, conf)
}

func gracefulShutdown(server *http.Server, conf *config.Conf) {
	if conf.Server.ShutdownSleepDuration > 0 {
		log.Printf("Waiting for %s until no new connections should come in anymore.", conf.Server.ShutdownSleepDuration.String())
		time.Sleep(conf.Server.ShutdownSleepDuration)
	}
	if conf.Server.DrainTimeout > 0 {
		log.Printf("Draining idle connections for up to %s.", conf.Server.DrainTimeout.String())
		draining.Store(true)
		drainStart := time.Now()
		conns := atomic.LoadInt64(&openConnectionsCount)
		for conns > 0 && time.Since(drainStart) < conf.Server.DrainTimeout {
			log.Printf("Waiting for connections to drain. Remaining connections: %d", conns)
			time.Sleep(1 * time.Second)
			conns = atomic.LoadInt64(&openConnectionsCount)
		}
		log.Printf("Finished draining connections. Remaining connections: %d", conns)
	}
	log.Printf("Shutting down server with a timeout of %s.", conf.Server.ShutdownTimeout.String())
	var cancelFunc context.CancelFunc
	ctx, cancelFunc := context.WithTimeout(context.Background(), conf.Server.ShutdownTimeout)
	defer cancelFunc()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error during server shutdown: %s", err)
	} else {
		log.Printf("Server shutdown completed successfully.")
	}
}
