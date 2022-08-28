package middlewares

import (
	"context"
	"net/http"

	"github.com/itsjoniur/bitlygo/internal/durable"
)

func Logger(logger *durable.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), 1, logger)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
