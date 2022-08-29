package middlewares

import (
	"context"
	"net/http"

	"github.com/itsjoniur/bitlygo/internal/durable"
)

// ContextMiddleware put durable.Database into context
func ContextMiddleware(db *durable.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), 10, db)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
