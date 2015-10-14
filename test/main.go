package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const path = "/src/github.com/Bredgren/gohtmlctrl/test"

var goPath = os.Getenv("GOPATH")

func main() {
	e := os.Chdir(filepath.Join(goPath, path))
	if e != nil {
		log.Fatal(e)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			http.ServeFile(w, r, "./"+r.RequestURI)
			return
		}
		http.ServeFile(w, r, "./index.html")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
