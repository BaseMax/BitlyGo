package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8000"
	fmt.Printf("Server is running on %v...\n", port)

	err := http.ListenAndServe(port, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s \n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
