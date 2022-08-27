package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/itsjoniur/bitlygo/internal/durable"
)

func ContextMiddleware(db *durable.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Printf("%#v", db)
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), 10, db)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
