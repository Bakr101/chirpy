package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(resWrite http.ResponseWriter, req *http.Request){
	type chirpReq struct{
		Body string `json:"body"`
	}
	type returnVals struct{
		Valid bool `json:"valid"`
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
		respondWithError(resWrite, http.StatusBadRequest,"Error encoding response body: %s", err)
		return	
	}	

	respondWithJSON(resWrite, http.StatusOK, returnVals{
		Valid: true,
	})
}