package main

import (
	"fmt"
	"net/http"
)

func main() {

	port := ":8080"
	protected := authMiddleware(adminHandler)
	http.HandleFunc("/start-job", protected)

	http.ListenAndServe(port, nil)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleBasicAuth(next, w, r)
	}
}

func handleBasicAuth(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()

	if !ok || user != "admin" || pass != "password123" {
		w.Header().Set("WWW-Authenticate", `Basic realm = Restricted`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	next(w, r)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "handled the auth")
}
