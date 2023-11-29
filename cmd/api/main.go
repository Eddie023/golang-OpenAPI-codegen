package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/internal/build"
	"github.com/eddie023/wex-tag/pkg/api"
	"github.com/eddie023/wex-tag/pkg/api/service"
	"github.com/eddie023/wex-tag/pkg/config"
	"github.com/eddie023/wex-tag/pkg/logger"
	"github.com/eddie023/wex-tag/pkg/types"

	_ "github.com/lib/pq"
)

func main() {
	l := logger.SlogWithColors()

	slog.Info("System starting...")

	if err := run(l); err != nil {
		slog.Error("startup", "error", err)
		os.Exit(1)
	}
}

func run(slog *slog.Logger) error {
	// GOMAXPROCS
	slog.Debug("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build.Build)

	cfg, err := config.GetParsedConfig()
	if err != nil {
		return err
	}

	slog.Debug("using config", "config", cfg)

	// try to connect to postgres server
	// fail early if unable to connect to postgres db
	db, err := ConnectDb(cfg)
	if err != nil {
		slog.Error("unable to connect to db")
		return err
	}

	swagger, err := types.GetSwagger()
	if err != nil {
		slog.Error("swagger spec", "err", err)
		os.Exit(1)
	}

	swagger.Servers = nil

	api := api.API{
		Config:  cfg,
		Db:      *db,
		Swagger: swagger,
		Logger:  slog,
		TransactionService: &service.Service{
			Ent: db,
		},
	}

	server := &http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: api.Handler(),
	}

	// server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrput/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				slog.Error("graceful shutdown failed", "deadline exceeded", true)
			}
		}()

		// Trigger graceful shutdown
		slog.Warn("gracefully shutting down server", "deadline exceeded", false)
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		serverStopCtx()
	}()

	go func() {
		slog.Info("server listening on", "port", cfg.Web.APIHost, "host", cfg.Web.APIHost)

		err = server.ListenAndServe()
		if err != nil {
			// Add sentry exception here
			log.Fatal(err)
		}
	}()

	<-serverCtx.Done()

	return nil
}

func ConnectDb(cfg *config.ApiConfig) (*ent.Client, error) {
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Dbname, cfg.Db.Password)

	slog.Info("connecting to db", "database", "postgres", "connection-url", connectionURL)

	client, err := ent.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}
