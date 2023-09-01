package main

import "net/http"

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, "OK")
}
