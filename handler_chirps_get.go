package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig)handlerGetChirps(resWrite http.ResponseWriter, req *http.Request){
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error getting chirps", err)
		return
	}
	var chirpsSlice []Chirp
	for _, chirp := range chirps{
		chirpJson := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_ID: chirp.UserID,
		}
		chirpsSlice = append(chirpsSlice, chirpJson)
	}
	respondWithJSON(resWrite, http.StatusOK, chirpsSlice)
}

func (cfg *apiConfig)handlerGetChirp(resWrite http.ResponseWriter, req *http.Request){
	id := req.PathValue("id")
	uuid, err:= uuid.Parse(id)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error parsing Id", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), uuid)
	if err != nil {
		respondWithError(resWrite, http.StatusNotFound, "ID Not Found", err)
		return
	}
	chirpJson := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_ID: chirp.UserID,
	}
	respondWithJSON(resWrite, http.StatusOK, chirpJson)
}