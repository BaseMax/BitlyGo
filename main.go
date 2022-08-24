package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id        uint   `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
	DeletedAt int64  `json:"deleted_at" db:"deleted_at"`
}

type Link struct {
	Id        uint   `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Url       string `json:"url" db:"url"`
	Visits    uint   `json:"visits" db:"visits"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
	DeletedAt int64  `json:"deleted_at" db:"deleted_at"`
}

type UserResponse struct {
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
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

func (u *User) Response(apiKey string) UserResponse {
	return UserResponse{
		Username: u.Username,
		APIKey:   apiKey,
	}
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

const (
	Host     = "localhost"
	Port     = 5432
	UserDB   = "postgres"
	Password = "1234"
	Database = "bitlygo"
)

var (
	FakeUserDB UsersRepo
	FakeLinkDB LinksRepo
	db         *pgxpool.Pool
)

func main() {
	port := ":8000"

	databaseUri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", UserDB, Password, Host, Port, Database)
	poolConfig, err := pgxpool.ParseConfig(databaseUri)
	CheckError(err)

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	CheckError(err)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(HeaderMiddleware)

	// Handlers
	r.Get("/", RootHandler)
	r.Post("/users", UserRegisterHandler)
	r.Get("/links", ShowLinksHandler)
	r.Post("/links/add", AddLinkHandler)
	r.Get("/links/top", TopLinksHandler)

	fmt.Printf("Server is running on %v...\n", port)

	err = http.ListenAndServe(port, r)
	CheckError(err)
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Documentation"))
}

func UserRegisterHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	var key string
	err := json.NewDecoder(req.Body).Decode(&user)
	CheckError(err)
	// save user into the database
	err = db.QueryRow(context.Background(), `insert into users(username, password) values ($1, $2) returning id`, user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "internal server error",
		})
	}

	err = db.QueryRow(context.Background(), `insert into api_keys(user_id, key) values ($1, $2) returning key`, user.Id, uuid.New()).Scan(&key)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Response(key))
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
	_, err = url.ParseRequestURI(link.Url)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "Invalid URL",
		})
		return
	}
	FakeLinkDB.Add(link)
	json.NewEncoder(w).Encode(link.Response())
}

func TopLinksHandler(w http.ResponseWriter, req *http.Request) {
	limitParam := req.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)
	CheckError(err)
	if limit < 1 || limit > 100 {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "limit value should be between 1-100",
		})
		return
	}
	links := FakeLinkDB.Links
	if limit > len(links) {
		limit = len(links)
	}
	sort.Slice(links, func(i, j int) bool {
		return links[i].Visits > links[j].Visits
	})
	json.NewEncoder(w).Encode(map[string]any{
		"status": true,
		"items":  links[:limit],
	})
}

func SelectUserByApiKey(apiKey string) {
	var id uint
	// 	var userFromDB UserModel
	err := db.QueryRow(context.Background(), `select id from api_keys where key = $1`, apiKey).Scan(&id)
	CheckError(err)
	fmt.Println(id)
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
