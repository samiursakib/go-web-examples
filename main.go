package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "requested path: %s", r.URL.Path)
	})

	r.HandleFunc("/books/{title}/{page}", func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Println("title:", vars["title"], "page", vars["page"])
		fmt.Fprintf(w, "title: %s, page: %s", vars["title"], vars["page"])
	})
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.ListenAndServe(":8080", r)
}
