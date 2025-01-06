package main

import (
	"net/http"

	"github.com/Bakr101/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerRevoke(resWrite http.ResponseWriter, req *http.Request){
	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(req.Context(), userToken)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	respondWithJSON(resWrite, http.StatusNoContent, nil)
}