package responses

import (
	"context"
	"net/http"

	"github.com/unrolled/render"

	"github.com/itsjoniur/bitlygo/internal/models"
)

type UserResponse struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
}

func (u UserResponse) Create(user *models.User, apiKey *models.ApiKey) *UserResponse {
	u.Username = user.Username
	u.ApiKey = apiKey.Key.String()
	return &u
}

func RenderUserResponse(ctx context.Context, w http.ResponseWriter, user *models.User, apiKey *models.ApiKey) {
	r := ctx.Value(2).(*render.Render)
	resp := UserResponse{}.Create(user, apiKey)

	r.JSON(w, http.StatusCreated, resp)
}
