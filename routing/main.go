package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"] // the book title slug
		page := vars["page"]   // the page

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	fmt.Println("Listening on port 80...")
	http.ListenAndServe(":80", r)
}
