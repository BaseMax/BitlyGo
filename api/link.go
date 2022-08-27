package api

import "net/http"

func addLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("add link"))
}

func addLinkByPathHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("add link by path"))
}

func updateLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("update link"))
}

func deleteLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("delete link"))
}

func searchLinkHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("search link"))
}

func showTopLinksHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("top links"))
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("redirect"))
}
