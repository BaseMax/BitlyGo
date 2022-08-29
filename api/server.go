package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/unrolled/render"

	"github.com/itsjoniur/bitlygo/internal/durable"
	"github.com/itsjoniur/bitlygo/internal/middlewares"
)

// StartAPI start an API on given port
func StartAPI(logger *durable.Logger, db *pgxpool.Pool, port string) error {
	router := chi.NewRouter()
	database := durable.WrapDatabase(db)

	// Setup middlewares
	router.Use(middlewares.Logger(logger)) //fs logger
	router.Use(middlewares.Header)
	router.Use(middlewares.ContextMiddleware(database))
	router.Use(middlewares.Render(render.New()))
	router.Use(middleware.Logger) // http requests logger
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	// Register routes
	router.Get("/", rootHandler)
	router.Post("/add", addLinkHandler)
	router.Post("/{name}", addLinkByPathHandler)
	router.Put("/{name}", updateLinkHandler)
	router.Delete("/{name}", deleteLinkHandler)
	router.Get("/{name}", redirectHandler)
	router.Get("/search", searchLinkHandler)
	router.Get("/top", showTopLinksHandler)
	router.Get("/expire-soon", showExpireSoonLinksHandler)

	// Serve HTTP
	log.Printf("Server running on %v port...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		return err
	}

	return nil
}
