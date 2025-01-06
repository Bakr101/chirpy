package main

import (
	"encoding/json"
	"net/http"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/google/uuid"
)

type PolkaReq struct{
	Event string `json:"event"`
	Data  map[string]string `json:"data"`
}

func (cfg *apiConfig)handlerPolka(resWrite http.ResponseWriter, req *http.Request){
	reqAPIKey, err := auth.GetAPIKey(req.Header)
	if err != nil || reqAPIKey != cfg.polkaKey{
		respondWithError(resWrite, http.StatusUnauthorized, "Invalid API Key", err)
		return
	}
	polkaParams := PolkaReq{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&polkaParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	if polkaParams.Event != "user.upgraded"{
		respondWithError(resWrite, http.StatusNoContent, "err", err)
		return
	}
	parsedUUID, err := uuid.Parse(polkaParams.Data["user_id"])
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "err parsing uuid", err)
		return
	}

	user, err := cfg.db.SetChirpyRed(req.Context(), parsedUUID)
	if err != nil {
		respondWithError(resWrite, http.StatusNotFound, "user not found", err)
		return
	}
	if !user.IsChirpyRed {
		respondWithError(resWrite, http.StatusNotFound, "couldn't set user to chirpy red", err)
		return
	}
	respondWithJSON(resWrite, http.StatusNoContent, nil)
}