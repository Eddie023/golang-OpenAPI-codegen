package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/eddie023/wex-tag/ent"
	"github.com/eddie023/wex-tag/pkg/config"
	"github.com/eddie023/wex-tag/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"

	"github.com/go-chi/httplog/v2"
	httpMiddleware "github.com/oapi-codegen/nethttp-middleware"
)

type API struct {
	Config             *config.ApiConfig
	Swagger            *openapi3.T
	Db                 ent.Client
	Logger             *slog.Logger
	TransactionService TransactionService
}

func (a *API) Handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(httplog.RequestLogger(getChiSlogLogger(a.Logger)))
	router.Use(httplog.RequestLogger(getChiSlogLogger(a.Logger)))
	router.Use(cors.Default().Handler)

	// healthcheck handler
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))

		// apiout.JSON(ctx, w, out, http.StatusCreated)
	})

	router.Group(func(r chi.Router) {
		r.Use(httpMiddleware.OapiRequestValidator(a.Swagger))
		types.HandlerWithOptions(a, types.ChiServerOptions{
			BaseRouter: r,
		})
	})

	return router
}

// getChiSlogLogger will initiate a structured logging for chi logger middleware.
func getChiSlogLogger(s *slog.Logger) *httplog.Logger {
	return &httplog.Logger{
		Logger: s,
		Options: httplog.Options{
			LogLevel: slog.LevelDebug,

			MessageFieldName: "message",
			LevelFieldName:   "severity",
			TimeFieldFormat:  time.RFC3339,

			QuietDownPeriod: 10 * time.Second,
		},
	}
}
