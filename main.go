package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

type UserResponse struct {
	Username string    `json:"username"`
	APIKey   uuid.UUID `json:"api_key"`
}

type UsersRepo struct {
	Users []User
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

func (u *UsersRepo) Add(user User) {
	u.Users = append(u.Users, user)
}

var (
	FakeUserDB UsersRepo
)

func main() {
	port := ":8000"
	// Handlers
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/users", UserRegisterHandler)

	fmt.Printf("Server is running on %v...\n", port)

	err := http.ListenAndServe(port, logRequest(http.DefaultServeMux))
	CheckError(err)
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Documentation"))
}

func UserRegisterHandler(w http.ResponseWriter, req *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&user)
	CheckError(err)
	user.Create()
	FakeUserDB.Add(user)
	json.NewEncoder(w).Encode(user.Response())
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s \n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
