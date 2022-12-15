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

	destination := r.FormValue("destination")

	_, err := net.DialTimeout("tcp", destination, time.Second*10)
	if err != nil {
		tmpl.Execute(w, struct {
			Success     bool
			Destination string
			Error       error
		}{false, destination, err})
		log.Printf("Error dialing %s: %s", destination, err)
		return
	}

	tmpl.Execute(w, struct {
		Success     bool
		Destination string
	}{true, destination})
	log.Printf("Successfully connected to %s", destination)
}

func main() {
	http.HandleFunc("/", testConnection)
	log.Print("Running on 8080")
	http.ListenAndServe(":8080", nil)
}
