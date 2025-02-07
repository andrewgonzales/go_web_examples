package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on port 80...")
	http.ListenAndServe(":80", nil)
}
