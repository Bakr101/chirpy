package main

import (
	"net/http"
	"time"

	"github.com/Bakr101/chirpy/internal/auth"
)

type RefreshRes struct{
	Token	string	`json:"token"`
}

func (cfg *apiConfig)handlerRefresh(resWrite http.ResponseWriter, req *http.Request){
	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	refreshTokenData, err := cfg.db.GetUserFromRefreshToken(req.Context(), userToken)
	if err != nil || time.Now().After(refreshTokenData.ExpiresAt) || refreshTokenData.RevokedAt.Valid{
		respondWithError(resWrite, http.StatusUnauthorized, "token expired", err)
		return
	}
	userAccessToken, err := auth.MakeJWT(refreshTokenData.UserID, cfg.tokenSecret, time.Duration(1))
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "error creating access token", err)
		return
	}
	respondWithJSON(resWrite, http.StatusOK, RefreshRes{
		Token: userAccessToken,
	})
}