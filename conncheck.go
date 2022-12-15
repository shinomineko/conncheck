package main

import (
	"embed"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

//go:embed tmpl/form.html
var formTmpl embed.FS

func testConnection(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(formTmpl, "tmpl/form.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	dest := r.FormValue("destination")

	_, err := net.DialTimeout("tcp", dest, time.Second*5)
	if err != nil {
		tmpl.Execute(w, struct {
			Success     bool
			Destination string
			Error       error
		}{false, dest, err})
		log.Printf("Error dialing %s: %s", dest, err)
		return
	}

	tmpl.Execute(w, struct {
		Success     bool
		Destination string
	}{true, dest})
	log.Printf("Successfully connected to %s", dest)
}

func main() {
	http.HandleFunc("/", testConnection)

	log.Print("Starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
