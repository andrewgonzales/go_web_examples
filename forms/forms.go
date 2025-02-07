package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("method %s\n", r.Method)
		if r.Method != http.MethodPost {
			// reset form for testing
			tmpl.Execute(w, nil)
			return
		}

		details := ContactDetails{
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		// do something with details
		fmt.Printf("details %#v\n", details)

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	fmt.Println("Listening on port 80...")
	http.ListenAndServe(":80", nil)
}
