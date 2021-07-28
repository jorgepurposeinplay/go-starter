package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/oakeshq/go-starter/config"
	"github.com/oakeshq/go-starter/internal/healthcheck"
	"github.com/oakeshq/go-starter/internal/storage"
	"github.com/oakeshq/go-starter/pkg/httperr"
	"github.com/oakeshq/go-starter/pkg/logs"
	gmiddleware "github.com/oakeshq/go-starter/pkg/middleware"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/oakeshq/go-starter/pkg/router"
)

// API exposes the integral struct
type API struct {
	handler http.Handler
	r       *router.Router
	config  *config.Config
	db      *gorm.DB
	service *storage.Service
}

// NewAPI instantiates a new REST API.
func NewAPI(
	config *config.Config,
	r *router.Router,
	db *gorm.DB,
) *API {
	api := &API{
		r:      r,
		config: config,
		db:     db,
	}
	repo := storage.NewRepository(db)
	service := storage.NewService(repo)
	api.service = service
	ctx := context.Background()
	r.Chi.Use(middleware.RealIP)
	r.Use(gmiddleware.RequestIDCtx)
	r.Use(httperr.Recoverer)

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	r.UseBypass(logs.NewStructuredLogger(logger))

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
	})

	healthcheck.RegisterHandlers(r)

	api.handler = corsHandler.Handler(chi.ServerBaseContext(ctx, r))
	return api
}
