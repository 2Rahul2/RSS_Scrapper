package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// (writes binary direct to http response)
	data, err := json.Marshal(payload) //converts json to string , Marshal returns data in bytes
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Failed to Marshal json response :%v", payload)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
