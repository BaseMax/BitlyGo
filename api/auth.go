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
		// return user exist
		return
	}

	user, err := models.CreateUser(req.Context(), authParam.Username, authParam.Password)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		// User Is exist
		w.Write([]byte("duplicated"))
		return
	}

	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	// render user response
	json.NewEncoder(w).Encode(user)
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
		// return user does not exist
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}
