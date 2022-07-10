package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", proxyToTwitter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func originAllowed(origin string) bool {
	for _, s := range allowedOrigins {
		if s == origin {
			return true
		}
	}
	for _, prefix := range allowedOriginPrefixes {
		if strings.HasPrefix(origin, prefix) {
			return true
		}
	}
	return false
}

func proxyToTwitter(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		handleCORSPreflight(w, req)
		return
	}
	ctx := context.Background()
	twReq := req.Clone(ctx)
	client := http.DefaultClient

	twReq.Host = ""
	twReq.URL.Host = "api.twitter.com"
	twReq.URL.Scheme = "https"
	twReq.URL.User = nil
	twReq.RequestURI = ""

	twResp, err := client.Do(twReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range twResp.Header {
		w.Header()[k] = v
	}
	origin := req.Header.Get("Origin")
	if originAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Expose-Headers", "x-rate-limit-limit, x-rate-limit-remaining, x-rate-limit-reset, x-response-time")

	w.WriteHeader(twResp.StatusCode)
	io.Copy(w, twResp.Body)
}

func handleCORSPreflight(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	if originAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.WriteHeader(http.StatusNoContent)
}
