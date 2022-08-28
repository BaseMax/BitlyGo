package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/itsjoniur/bitlygo/internal/durable"
	"github.com/itsjoniur/bitlygo/internal/middlewares"
	"github.com/jackc/pgx/v4/pgxpool"
)

func StartAPI(logger *log.Logger, db *pgxpool.Pool, port string) error {
	router := chi.NewRouter()
	database := durable.WrapDatabase(db)
	// setup middlewares
	router.Use(middlewares.Header)
	router.Use(middlewares.ContextMiddleware(database))
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	// register routes
	router.Get("/", rootHandler)
	router.Post("/add", addLinkHandler)
	router.Post("/{name}", addLinkByPathHandler)
	router.Put("/{name}", updateLinkHandler)
	router.Delete("/{name}", deleteLinkHandler)
	router.Get("/{name}", redirectHandler)
	router.Get("/search", searchLinkHandler)
	router.Get("/top", showTopLinksHandler)

	log.Printf("Server running on %v port...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		return err
	}

	return nil
}
