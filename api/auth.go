package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/itsjoniur/bitlygo/internal/models"
	"github.com/itsjoniur/bitlygo/internal/responses"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserRegisterHandler(w http.ResponseWriter, req *http.Request) {
	authParam := Auth{}

	json.NewDecoder(req.Body).Decode(&authParam)

	if authParam.Username == "" {
		responses.FieldEmptyError(req.Context(), w, "username")
		return
	}

	if authParam.Password == "" {
		responses.FieldEmptyError(req.Context(), w, "password")
		return
	}

	isExist := models.GetUserByUsername(req.Context(), authParam.Username)
	if isExist != nil {
		responses.UserIsExistsError(req.Context(), w)
		return
	}

	user, err := models.CreateUser(req.Context(), authParam.Username, authParam.Password)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		responses.UserIsExistsError(req.Context(), w)
		return
	}

	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	key, err := models.CreateApiKey(req.Context(), user)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}
	// render user response
	responses.RenderUserResponse(req.Context(), w, user, key)
}

func UserLoginHandler(w http.ResponseWriter, req *http.Request) {
	authParam := Auth{}

	json.NewDecoder(req.Body).Decode(&authParam)

	if authParam.Username == "" {
		responses.FieldEmptyError(req.Context(), w, "username")
		return
	}

	if authParam.Password == "" {
		responses.FieldEmptyError(req.Context(), w, "password")
		return
	}

	user := models.GetUserByUsername(req.Context(), authParam.Username)
	if user == nil {
		responses.NotFoundError(req.Context(), w)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}
