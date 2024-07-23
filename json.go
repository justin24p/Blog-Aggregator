package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// returns as bytes
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal JSON %v", payload)
		w.WriteHeader(500)
		return 
	}
	w.Header().Add("Content-Type","applicatoin/json")
	w.WriteHeader(200)
	w.Write(data)
}