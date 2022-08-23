package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8000"
	fmt.Printf("Server is running on %v...\n", port)

	http.ListenAndServe(port, nil)
}
