package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/itsjoniur/bitlygo/internal/models"
)

func addLinkHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}

	params := Params{}
	json.NewDecoder(req.Body).Decode(&params)

	if params.Name == "" {
		// Generate random string
		params.Name = "aaa" // This will be the default value until we implement a function to generate random string for us
	}

	if params.Link == "" {
		// Link is a required field and when it's empty we should return an error
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]any{
			"status":  false,
			"message": "link can not be ampty",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if _, err := url.ParseRequestURI(params.Link); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]any{
			"status":  false,
			"message": "link must be a vaild url",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	link, err := models.CreateLink(req.Context(), 0, params.Name, params.Link)
	if err != nil && strings.Contains(string(err.Error()), "duplicate key") {
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]any{
			"status":  false,
			"message": fmt.Sprintf("link with name `%v` exists", params.Name),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]any{
			"status":  false,
			"message": http.StatusText(http.StatusInternalServerError),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(link)
}

func addLinkByPathHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("add link by path"))
}

func updateLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("update link"))
}

func deleteLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("delete link"))
}

func searchLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("search link"))
}

func showTopLinksHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("top links"))
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("redirect"))
}
