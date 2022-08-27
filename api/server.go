package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartAPI(logger *log.Logger, port string) error {
	router := chi.NewRouter()
	// setup middlewares
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

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		return err
	}

	return nil
}
