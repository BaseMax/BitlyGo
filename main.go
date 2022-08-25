package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id        uint       `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	Password  string     `json:"password" db:"password"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Link struct {
	Id            uint       `json:"id" db:"id"`
	OwnerId       *uint      `json:"owner_id" db:"owner_id"`
	Name          string     `json:"name" db:"name"`
	Url           string     `json:"url" db:"url"`
	Visits        uint       `json:"visits" db:"visits"`
	StatisticsKey *string    `json:"statistics_key" db:"statistics_key"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	ExpiredAt     *time.Time `json:"expired_at" db:"expired_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
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

type LinkUpdate struct {
	NewName string `json:"new_name"`
	Url     string `json:"url"`
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

const (
	Host     = "localhost"
	Port     = 5432
	UserDB   = "postgres"
	Password = "1234"
	Database = "bitlygo"
)

var (
	db          *pgxpool.Pool
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	r.Get("/links/search", SearchLinkHandler)
	r.Get("/links/top", TopLinksHandler)
	r.Put("/links/{name}", UpdateLinkHandler)
	r.Get("/{name}", RedirectHandler)

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

func RedirectHandler(w http.ResponseWriter, req *http.Request) {
	urlName := chi.URLParam(req, "name")
	var url string
	err := db.QueryRow(context.Background(), `select link from links where name = $1`, urlName).Scan(&url)
	CheckError(err)
	_, err = db.Exec(context.Background(), `update links set visits = coalesce(visits, 0) + 1 where name = $1`, urlName)
	CheckError(err)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

func ShowLinksHandler(w http.ResponseWriter, req *http.Request) {
	res := map[string]any{
		"status": true,
		"items":  nil,
	}
	json.NewEncoder(w).Encode(res)
}

func AddLinkHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := req.Header.Get("API-KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "api key is required",
		})
	}
	var link Link
	err := json.NewDecoder(req.Body).Decode(&link)
	CheckError(err)
	_, err = url.ParseRequestURI(link.Url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "Invalid URL",
		})
		return
	}
	if link.Name == "" {
		// generate random key
		link.Name = RandStringRunes(6)
	}
	user := GetUserByApiKey(apiKey)
	_, err = db.Exec(context.Background(), `insert into links(owner_id, name, link) values($1, $2, $3)`, user.Id, link.Name, link.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "something went wrong please try again",
		})
		return
	}
	link.OwnerId = &user.Id
	json.NewEncoder(w).Encode(link.Response())

}

func TopLinksHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := req.Header.Get("API-KEY")
	limitParam := req.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "10"
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "invalid value for limit parameter",
		})
		return
	}
	if limit < 1 || limit > 100 {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "limit value should be between 1-100",
		})
		return
	}
	if apiKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "missing api key",
		})
		return
	}
	user := GetUserByApiKey(apiKey)
	links := []LinkResponse{}
	rows, _ := db.Query(context.Background(), `select * from links where owner_id = $1 order by visits desc limit $2`, user.Id, limit)
	for rows.Next() {
		link := &Link{}
		err := rows.Scan(&link.Id, &link.OwnerId, &link.Name, &link.Url, &link.Visits, &link.CreatedAt, &link.UpdatedAt, &link.DeletedAt, &link.ExpiredAt, &link.StatisticsKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "something went wrong please try again",
			})
		}
		links = append(links, link.Response())
	}
	json.NewEncoder(w).Encode(map[string]any{
		"status": true,
		"items":  links,
	})
}

func SearchLinkHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := req.Header.Get("API-KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "API key is required",
		})
		return
	}

	searchQuery := strings.Trim(req.URL.Query().Get("q"), "\\r\\n")
	limitParam := req.URL.Query().Get("limit")
	if searchQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "query parameter is required",
		})
		return
	}
	if limitParam == "" {
		limitParam = "10"
	}
	limit, err := strconv.Atoi(limitParam)
	CheckError(err)

	result := make(map[string]any)
	result["items"] = make(map[string]string)
	user := GetUserByApiKey(apiKey)
	dbQuery := fmt.Sprintf("select name, link from links where owner_id = $1 and name like '%%%v%%' limit $2", searchQuery)
	rows, _ := db.Query(context.Background(), dbQuery, user.Id, limit)
	for rows.Next() {
		var name string
		var link string
		err := rows.Scan(&name, &link)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"status":  false,
				"message": "something went wrong please try again",
			})
			return
		}
		result["items"].(map[string]string)[name] = link
	}
	result["status"] = true
	json.NewEncoder(w).Encode(result)
}

func UpdateLinkHandler(w http.ResponseWriter, req *http.Request) {
	apiKey := req.Header.Get("API-KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "api key is required",
		})
		return
	}
	user := GetUserByApiKey(apiKey)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "user does not exist",
		})
		return
	}
	var isExist bool
	urlName := chi.URLParam(req, "name")
	err := db.QueryRow(context.Background(), `select exists(select id from links where name = $1 and owner_id = $2)`, urlName, user.Id).Scan(&isExist)
	CheckError(err)
	if !isExist {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": fmt.Sprintf("link with name %v does not exist", urlName),
		})
		return
	}
	var link LinkUpdate
	json.NewDecoder(req.Body).Decode(&link)
	if link.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  false,
			"message": "link is required field",
		})
		return
	}

	dbQuery := "update links set link=$1"
	if link.NewName != "" {
		dbQuery += ", name=$2 where owner_id = $3 and name = $4"
		_, err := db.Exec(context.Background(), dbQuery, link.Url, link.NewName, user.Id, urlName)
		CheckError(err)
	} else {
		dbQuery += " where owner_id = $2 and name = $3"
		_, err := db.Exec(context.Background(), dbQuery, link.Url, user.Id, urlName)
		CheckError(err)
	}

	json.NewEncoder(w).Encode(map[string]bool{"status": true})

}

func GetUserByApiKey(apiKey string) *User {
	var id uint
	user := &User{}
	// 	var userFromDB UserModel
	err := db.QueryRow(context.Background(), `select user_id from api_keys where key = $1`, apiKey).Scan(&id)
	err = db.QueryRow(context.Background(), `select * from users where id = $1`, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil
	}
	return user
}

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	})
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
