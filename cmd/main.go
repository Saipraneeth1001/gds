package main

import (
	"fmt"
	"net/http"
)

func main() {

	port := ":8080"
	http.HandleFunc("/start-job", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "job initiated")
	})

	http.ListenAndServe(port, nil)
}
