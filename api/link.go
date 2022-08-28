package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/itsjoniur/bitlygo/internal/models"
	"github.com/itsjoniur/bitlygo/internal/responses"
	"github.com/itsjoniur/bitlygo/pkg/strutil"
)

func addLinkHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}
	var err error
	link := &models.Link{}
	apiKey := req.Header.Get("API-KEY")

	params := Params{}
	json.NewDecoder(req.Body).Decode(&params)

	params.Name, err = strutil.RemoveNonAlphanumerical(params.Name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if params.Name == "" {
		// Generate random string
		params.Name = strutil.RandStringRunes(8)
	}

	if params.Link == "" {
		// Link is a required field and when it's empty we should return an error
		responses.BadRequestError(req.Context(), w)
		return
	}

	if _, err := url.ParseRequestURI(params.Link); err != nil {
		responses.InvalidLinkError(req.Context(), w)
		return
	}

	if apiKey != "" {
		link, err = models.CreateLink(req.Context(), 0, params.Name, params.Link)
	} else {
		link, err = models.CreateLinkWithExpireTime(req.Context(), 0, params.Name, params.Link)
	}
	if err != nil && strings.Contains(string(err.Error()), "duplicate key") {
		responses.LinkIsExistsError(req.Context(), w, params.Name)
		return
	}

	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderNewLinkResponse(req.Context(), w, link)
}

func addLinkByPathHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}
	var err error
	link := &models.Link{}
	params := Params{}
	apiKey := req.Header.Get("API-KEY")

	json.NewDecoder(req.Body).Decode(&params)
	if params.Link == "" {
		responses.FieldEmptyError(req.Context(), w, "link")
		return
	}

	params.Name = chi.URLParam(req, "name")
	params.Name, err = strutil.RemoveNonAlphanumerical(params.Name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if params.Name == "" {
		params.Name = strutil.RandStringRunes(8)
	}

	if _, err := url.ParseRequestURI(params.Link); err != nil {
		responses.InvalidLinkError(req.Context(), w)
		return
	}

	if apiKey != "" {
		link, err = models.CreateLink(req.Context(), 0, params.Name, params.Link)
	} else {
		link, err = models.CreateLinkWithExpireTime(req.Context(), 0, params.Name, params.Link)
	}
	if err != nil && strings.Contains(string(err.Error()), "duplicate key") {
		responses.LinkIsExistsError(req.Context(), w, params.Name)
		return
	}

	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderNewLinkResponse(req.Context(), w, link)
}

func updateLinkHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		NewName string `json:"new_name"`
		Link    string `json:"link"`
	}
	var err error
	params := Params{}
	name := chi.URLParam(req, "name")

	name, err = strutil.RemoveNonAlphanumerical(name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	isExist := models.GetLinkByName(req.Context(), name)
	if isExist == nil {
		responses.NotFoundError(req.Context(), w)
		return
	}

	json.NewDecoder(req.Body).Decode(&params)
	if params.Link == "" {
		responses.FieldEmptyError(req.Context(), w, "link")
		return
	}

	link, err := models.UpdateLinkByName(req.Context(), name, params.NewName, params.Link)
	if err != nil && strings.Contains(string(err.Error()), "duplicate key") {
		responses.LinkIsExistsError(req.Context(), w, params.NewName)
		return
	}
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderNewLinkResponse(req.Context(), w, link)
}

func deleteLinkHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	name := chi.URLParam(req, "name")

	name, err = strutil.RemoveNonAlphanumerical(name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if name == "" {
		responses.BadRequestError(req.Context(), w)
		return
	}

	err = models.DeleteLinkByName(req.Context(), name)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func searchLinkHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	sq := req.URL.Query().Get("q")
	limit := req.URL.Query().Get("limit")

	sq, err = strutil.RemoveNonAlphanumerical(sq)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if len(sq) < 1 {
		responses.FieldEmptyError(req.Context(), w, "search")
		return
	}

	if limit == "" {
		limit = "10"
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	if 1 > l || l > 100 {
		responses.LimitRangeError(req.Context(), w)
		return
	}

	links, err := models.SearchLinkByName(req.Context(), sq, l)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderSearchLinkResponse(req.Context(), w, links)

}

func showTopLinksHandler(w http.ResponseWriter, req *http.Request) {
	limit := req.URL.Query().Get("limit")

	if limit == "" {
		limit = "10"
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	if 1 > l || 1 > 100 {
		responses.LimitRangeError(req.Context(), w)
		return
	}

	tl, err := models.TopLinksByVisits(req.Context(), l)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderTopLinksResponse(req.Context(), w, tl)
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	name := chi.URLParam(req, "name")

	name, err = strutil.RemoveNonAlphanumerical(name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if name == "" {
		responses.BadRequestError(req.Context(), w)
		return
	}

	link := models.GetLinkByName(req.Context(), name)
	if link == nil {
		responses.NotFoundError(req.Context(), w)
		return
	}

	go models.AddViewToLinkByName(req.Context(), name)

	http.Redirect(w, req, link.Link, http.StatusPermanentRedirect)
}

func showExpireSoonLinksHandler(w http.ResponseWriter, req *http.Request) {
	links, err := models.GetExpireSoonLinks(req.Context())
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderExpireLinkResponse(req.Context(), w, links)
}
