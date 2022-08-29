package api

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/itsjoniur/bitlygo/internal/configs"
	"github.com/itsjoniur/bitlygo/internal/responses"
)

// rootHandler show the project documentaion as a HTML page
func rootHandler(w http.ResponseWriter, req *http.Request) {
	dir, err := configs.GetRootDir()
	if err != nil {
		responses.InternalServerError(req.Context(), w)
	}
	dat, err := os.ReadFile(path.Join(dir, "README.md"))
	if err != nil {
		fmt.Println(err)
		responses.InternalServerError(req.Context(), w)
	}

	md := []byte(dat)
	docs := markdown.ToHTML(md, nil, nil)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(docs)
}
