package middlewares

import (
	"net/http"

	"github.com/itsjoniur/bitlygo/internal/models"
)

func CurrentUser(req *http.Request) *models.User {
	apiKey := req.Header.Get("API-KEY")

	if apiKey == "" {
		return nil
	}

	user := models.GetUserByApiKey(req.Context(), apiKey)
	return user
}
