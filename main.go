package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type User struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type Link struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	Visits    uint   `json:"visits"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type UserResponse struct {
	Username string    `json:"username"`
	APIKey   uuid.UUID `json:"api_key"`
}

type LinkResponse struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Visits uint   `json:"visits"`
}

type UsersRepo struct {
	Users []User
}

type LinksRepo struct {
	Links []Link
}

func (u *User) Create() {
	u.CreatedAt = time.Now().UTC().String()
	u.UpdatedAt = time.Now().UTC().String()
}

func (u *User) Response() UserResponse {
	key := uuid.New()
	return UserResponse{
		Username: u.Username,
		APIKey:   key,
	}
}

func (l *Link) Create() {
	l.CreatedAt = time.Now().UTC().String()
	l.UpdatedAt = time.Now().UTC().String()
}

func (l *Link) Response() LinkResponse {
	return LinkResponse{
		Name:   l.Name,
		Url:    l.Url,
		Visits: l.Visits,
	}
}

func (u *UsersRepo) Add(user User) {
	u.Users = append(u.Users, user)
}

func (l *LinksRepo) Add(link Link) {
	l.Links = append(l.Links, link)
}

var (
	FakeUserDB UsersRepo
	FakeLinkDB LinksRepo
)

func main() {
	port := ":8000"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(HeaderMiddleware)

	// Handlers
	r.Get("/", RootHandler)
	r.Post("/users", UserRegisterHandler)
	r.Get("/links", ShowLinksHandler)
	r.Post("/links/add", AddLinkHandler)

	fmt.Printf("Server is running on %v...\n", port)

	err := http.ListenAndServe(port, r)
	CheckError(err)
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Documentation"))
}

func UserRegisterHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	CheckError(err)
	user.Create()
	FakeUserDB.Add(user)
	json.NewEncoder(w).Encode(user.Response())
}

func ShowLinksHandler(w http.ResponseWriter, req *http.Request) {
	res := map[string]any{
		"status": true,
		"items":  FakeLinkDB.Links,
	}
	json.NewEncoder(w).Encode(res)
}

func AddLinkHandler(w http.ResponseWriter, req *http.Request) {
	var link Link
	err := json.NewDecoder(req.Body).Decode(&link)
	CheckError(err)
	link.Create()
	FakeLinkDB.Add(link)
	json.NewEncoder(w).Encode(link.Response())
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
