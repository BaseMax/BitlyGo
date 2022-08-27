package api

import "net/http"

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("documentation"))
}
