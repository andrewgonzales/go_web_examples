package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Handle dynamic params
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("params %v \n", r.URL.Query())
		fmt.Fprint(w, "Welcome to my website!")
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on port 80...")

	http.ListenAndServe(":80", nil)
}
