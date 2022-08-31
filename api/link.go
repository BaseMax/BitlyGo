package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"

	"github.com/itsjoniur/bitlygo/internal/middlewares"
	"github.com/itsjoniur/bitlygo/internal/models"
	"github.com/itsjoniur/bitlygo/internal/responses"
	"github.com/itsjoniur/bitlygo/pkg/strutil"
)

func createLink(w http.ResponseWriter, req *http.Request, name, link string) {
	var err error
	newLink := &models.Link{}

	name, err = strutil.RemoveNonAlphanumerical(name)
	if err != nil {
		responses.BadRequestError(req.Context(), w)
		return
	}

	if name == "" {
		name = models.GetUniqueName(req.Context(), 8)
	}

	if slices.Contains(ReservedNames(), name) {
		responses.ReservedNameError(req.Context(), w)
		return
	}

	if 4 > len(name) || len(name) > 25 {
		responses.LinkNameLengthError(req.Context(), w)
		return
	}

	if link == "" {
		// Link is a required field and when it's empty we should return an error
		responses.BadRequestError(req.Context(), w)
		return
	}

	if _, err := url.ParseRequestURI(link); err != nil {
		responses.InvalidLinkError(req.Context(), w)
		return
	}

	user := middlewares.CurrentUser(req)
	if user != nil {
		newLink, err = models.CreateLink(req.Context(), user.ID, name, link)
	} else {
		newLink, err = models.CreateLinkWithExpireTime(req.Context(), 0, name, link)
	}
	if err != nil && strings.Contains(string(err.Error()), "duplicate key") {
		responses.LinkIsExistsError(req.Context(), w, name)
		return
	}

	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderNewLinkResponse(req.Context(), w, newLink)
}

// AddLinkHandler store new links in database
func addLinkHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}
	params := Params{}

	json.NewDecoder(req.Body).Decode(&params)
	createLink(w, req, params.Name, params.Link)

}

// AddLinkByPathHandler same addLinkHandler but get the link name from url path
func addLinkByPathHandler(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Name string `json:"name"`
		Link string `json:"link"`
	}
	params := Params{}

	json.NewDecoder(req.Body).Decode(&params)
	createLink(w, req, params.Name, params.Link)
}

// UpdateLinkHandler update link and name values of a Link
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

	if slices.Contains(ReservedNames(), params.NewName) {
		responses.ReservedNameError(req.Context(), w)
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

// DeleteLinkHandler delete a link from database
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

	if slices.Contains(ReservedNames(), name) {
		responses.ReservedNameError(req.Context(), w)
		return
	}

	err = models.DeleteLinkByName(req.Context(), name)
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

// SearchLinkHandler find matched links
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

	var links []*models.Link
	user := middlewares.CurrentUser(req)
	if user != nil {
		links, err = models.SearchLinkByName(req.Context(), sq, user.ID, l)
	} else {
		links, err = models.SearchLinkByName(req.Context(), sq, 0, l)
	}
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderSearchLinkResponse(req.Context(), w, links)

}

// ShowTopLinksHandler show top links by visits
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

	var tl []*models.Link
	user := middlewares.CurrentUser(req)
	if user != nil {
		tl, err = models.TopLinksByVisits(req.Context(), user.ID, l)
	} else {
		tl, err = models.TopLinksByVisits(req.Context(), 0, l)
	}
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderTopLinksResponse(req.Context(), w, tl)
}

// RedirectHandler redirect to target link
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

// ShowTopLinksHandler show links will be expired soon
func showExpireSoonLinksHandler(w http.ResponseWriter, req *http.Request) {
	var links []*models.Link
	var err error

	user := middlewares.CurrentUser(req)
	if user != nil {
		links, err = models.GetExpireSoonLinks(req.Context(), user.ID)
	} else {
		links, err = models.GetExpireSoonLinks(req.Context(), 0)
	}
	if err != nil {
		responses.InternalServerError(req.Context(), w)
		return
	}

	responses.RenderExpireLinkResponse(req.Context(), w, links)
}

// ReservedNames return reserved names
func ReservedNames() []string {
	return []string{"add", "top", "search", "expire-soon"}
}
