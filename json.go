package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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