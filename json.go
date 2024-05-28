package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//Create two functions:
//
//respondWithJSON(w http.ResponseWriter, status int, payload interface{})
//respondWithError(w http.ResponseWriter, code int, msg string) (which calls respondWithJSON with error-specific values

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println(
			"Responding with 5xx error: %v",
			msg,
		)
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(
		w,
		code,
		errResponse{Error: msg},
	)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	log.Println(
		"Failed to marshal JSON response %v",
		payload,
	)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Add(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(code)
	w.Write(dat)
}
