package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(resWrite http.ResponseWriter, req *http.Request){
	type chirpReq struct{
		Body string `json:"body"`
	}
	type returnVals struct{
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := chirpReq{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	const maxChirpLength = 140
	if len(reqParams.Body) > maxChirpLength {
		respondWithError(resWrite, http.StatusBadRequest,"Chirp is too long", nil)
		return	
	}

	respondWithJSON(resWrite, http.StatusOK, returnVals{
		Body: handleBadwords(reqParams.Body),
	})
}

func handleBadwords(str string) string{
	splitted := strings.Split(str, " ")
	lowered := strings.ToLower(str)
	splittedToClean := strings.Split(lowered, " ")
	cleaned := []string{}
	badWords := map[string]string{
		"kerfuffle": "kerfuffle",
		"sharbert": "sharbert",
		"fornax": "fornax",
	}

	for idx, word := range splittedToClean{
		_, exists := badWords[word]
		if exists {
			cleaned = append(cleaned, "****")
		}
		if !exists {
			cleaned = append(cleaned, splitted[idx])
		}
	}

	return strings.Join(cleaned, " ")
}