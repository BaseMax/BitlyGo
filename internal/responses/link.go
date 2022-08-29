package responses

import (
	"context"
	"net/http"
	"time"

	"github.com/itsjoniur/bitlygo/internal/models"
	"github.com/unrolled/render"
)

// Struct for a link response
type LinkResponse struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	Visits int    `json:"visits"`
}

// Struct for new link response
type NewLinkResponse struct {
	Status bool `json:"status" default:"true"`
	Item   struct {
		Name      string     `json:"name"`
		Link      string     `json:"link"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		ExpiredAt *time.Time `json:"expired_at" default:"null"`
	} `json:"item"`
}

// Struct for top links response
type TopLinksResponse struct {
	Status bool           `json:"status" bool:"true"`
	Items  []LinkResponse `json:"items"`
}

// Struct for search result response
type SearchLinkResponse struct {
	Status bool              `json:"status" bool:"true"`
	Items  map[string]string `json:"items"`
}

// Struct for links will be expired soon response
type ExpireLinksResponse struct {
	Status bool              `json:"status" bool:"true"`
	Items  map[string]string `json:"items"`
}

// LinkResponse.Create create a link response
func (l LinkResponse) Create(ml *models.Link) *LinkResponse {
	l.Name = ml.Name
	l.Link = ml.Link
	l.Visits = ml.Visits
	return &l
}

// NewLinkResponse.Create create a new link response
func (l NewLinkResponse) Create(ml *models.Link) *NewLinkResponse {
	l.Status = true
	l.Item.Name = ml.Name
	l.Item.Link = ml.Link
	l.Item.CreatedAt = ml.CreatedAt
	l.Item.UpdatedAt = ml.UpdatedAt
	l.Item.ExpiredAt = ml.ExpiredAt
	return &l
}

// TopLinksResponse.Create create top links response
func (l TopLinksResponse) Create(items []LinkResponse) *TopLinksResponse {
	l.Status = true
	l.Items = items
	return &l
}

// SearchLinkResponse.Create create search link response
func (l SearchLinkResponse) Create(items map[string]string) *SearchLinkResponse {
	l.Status = true
	l.Items = items
	return &l
}

// ExpireLinksResponse.Create create expire soon links response
func (l ExpireLinksResponse) Create(items map[string]string) *ExpireLinksResponse {
	l.Status = true
	l.Items = items
	return &l
}

// RenderNewLinkResponse render a new link response
func RenderNewLinkResponse(ctx context.Context, w http.ResponseWriter, link *models.Link) {
	r := ctx.Value(2).(*render.Render)
	resp := NewLinkResponse{}.Create(link)

	r.JSON(w, http.StatusCreated, resp)
}

// RenderTopLinksResponse render a top links response
func RenderTopLinksResponse(ctx context.Context, w http.ResponseWriter, links []*models.Link) {
	r := ctx.Value(2).(*render.Render)
	linksResponse := []LinkResponse{}

	for _, link := range links {
		linksResponse = append(linksResponse, *LinkResponse{}.Create(link))
	}

	resp := TopLinksResponse{}.Create(linksResponse)
	r.JSON(w, http.StatusOK, resp)
}

// RenderSearchLinkResponse render a search link response
func RenderSearchLinkResponse(ctx context.Context, w http.ResponseWriter, links []*models.Link) {
	r := ctx.Value(2).(*render.Render)
	linksMap := map[string]string{}

	for _, link := range links {
		linksMap[link.Name] = link.Link
	}

	resp := SearchLinkResponse{}.Create(linksMap)
	r.JSON(w, http.StatusOK, resp)
}

// RenderExpireLinkResponse render an expire link response
func RenderExpireLinkResponse(ctx context.Context, w http.ResponseWriter, links []*models.Link) {
	r := ctx.Value(2).(*render.Render)
	linksMap := map[string]string{}

	for _, link := range links {
		linksMap[link.Name] = link.Link
	}

	resp := ExpireLinksResponse{}.Create(linksMap)
	r.JSON(w, http.StatusOK, resp)
}
