package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// formats json object
func respondWithError(w http.ResponseWriter, code int, msg string) {
	// 500 are server errors
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse {
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// returns as json encoded bytes 
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON %v", payload)
		w.WriteHeader(500)
		return 
	}
	w.Header().Add("Content-Type","applicatoin/json")
	w.WriteHeader(code)
	w.Write(data)
}