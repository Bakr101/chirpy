package main

import "net/http"

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