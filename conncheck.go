package main

import (
	"embed"
	"html"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

//go:embed tmpl/form.html
var formTmpl embed.FS

var defaultTimeoutSeconds int = 5

func testConnection(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(formTmpl, "tmpl/form.html"))

	nodeName := os.Getenv("NODE_NAME")
	if nodeName == "" {
		nodeName, _ = os.Hostname()
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, struct {
			DefaultTimeoutSeconds int
			NodeName              string
		}{defaultTimeoutSeconds, nodeName})
		return
	}

	dest := r.FormValue("destination")
	escapedDest := strings.Replace(dest, "\n", "", -1)
	escapedDest = strings.Replace(escapedDest, "\r", "", -1)

	_, err := net.DialTimeout("tcp", escapedDest, time.Second*time.Duration(defaultTimeoutSeconds))
	if err != nil {
		tmpl.Execute(w, struct {
			DefaultTimeoutSeconds int
			NodeName              string
			Success               bool
			Destination           string
			Error                 error
		}{defaultTimeoutSeconds, nodeName, false, html.EscapeString(dest), err})
		logger.Error("Error dialing", "dest", escapedDest, "error", err)
		return
	}

	tmpl.Execute(w, struct {
		DefaultTimeoutSeconds int
		NodeName              string
		Success               bool
		Destination           string
	}{defaultTimeoutSeconds, nodeName, true, html.EscapeString(dest)})
	logger.Info("Successfully connected", "dest", escapedDest)
}

func main() {
	http.HandleFunc("/", testConnection)

	logger.Info("Starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("Error starting server", "error", err)
	}
}
