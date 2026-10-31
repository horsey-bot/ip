package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type response struct {
	IP string `json:"ip"`
}

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>IP address</title>
</head>
<body>
  <p>Your IP address is:</p>
  <pre>{{.}}</pre>
</body>
</html>
`))

func clientIP(r *http.Request) string {
	// Caddy is the only public entry point. It overwrites this header with the
	// client address before forwarding the request.
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

func noStore(w http.ResponseWriter, contentType string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", contentType)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		noStore(w, "text/html; charset=utf-8")
		_ = page.Execute(w, clientIP(r))
	})

	mux.HandleFunc("GET /text", func(w http.ResponseWriter, r *http.Request) {
		noStore(w, "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(clientIP(r) + "\n"))
	})

	mux.HandleFunc("GET /json", func(w http.ResponseWriter, r *http.Request) {
		noStore(w, "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(response{IP: clientIP(r)})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
