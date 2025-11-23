package main

import (
	"net/http"

	"config/server"
	"database"

	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"time"

	"opentelemetry"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelSdkSetup := opentelemetry.NewOtelSdkSetup(opentelemetry.SetupOTelSDKOptions{
		ServiceName:    "config",
		ServiceVersion: "0.1.0",
	})
	// Set up OpenTelemetry.
	otelShutdown, err := otelSdkSetup.Setup(ctx)
	if err != nil {
		return err
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Set up router, database and app server.
	r := chi.NewRouter()
	db := database.Connect()
	db.AutoMigrate(&database.Device{}, &database.TestSession{}, &database.Scenario{}, &database.ScenarioCondition{}, &database.ConditionValue{})

	appSrv := server.NewServer(db, r)
	appSrv.Start()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":3002",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      otelhttp.NewHandler(r, "/"),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	InitMetrics()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return err
}
