package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(resWrite http.ResponseWriter, code int, msg string, err error){
	if err != nil {
		log.Println(err)
	}
	
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(resWrite, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(resWrite http.ResponseWriter, code int, payload interface{}){
	resWrite.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshaling Json: %s", err)
		resWrite.WriteHeader(500)
		return
	}
	resWrite.WriteHeader(code)
	resWrite.Write(dat)
}