package main

import (
	"net/http"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/Bakr101/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig)handlerDeleteChirp(resWrite http.ResponseWriter, req *http.Request){
	id := req.PathValue("id")
	chirpUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "Invalid UUID", err)
	}
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	userUUID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWrite, http.StatusForbidden, "Invalid user", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpUUID)
	if err != nil {
		respondWithError(resWrite, http.StatusNotFound, "chirp not found", err)
		return
	}
	if userUUID != chirp.UserID {
		respondWithError(resWrite, http.StatusForbidden, "unauthorized action", err)
		return
	}
	dbParams := database.DeleteChirpParams{
		ID: chirpUUID,
		UserID: userUUID,
	}
	_, err = cfg.db.DeleteChirp(req.Context(), dbParams)
	if err != nil {
		respondWithError(resWrite, http.StatusNotFound, "chirp not found", err)
	}
	respondWithJSON(resWrite, http.StatusNoContent, nil)
}