package main

import (
	"encoding/json"
	"net/http"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/Bakr101/chirpy/internal/database"
)



func (cfg *apiConfig)handlerChangeUserInfo(resWrite http.ResponseWriter, req *http.Request){
	reqParams := UserReq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
	}
	
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	userUUID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "Invalid user", err)
		return
	}
	user, err := cfg.db.GetUserByID(req.Context(), userUUID)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "Invalid user", err)
		return
	}
	newHashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error hashing password:", err)
		return
	}
	dbParams := database.ChangePassEmailParams{
		Email: reqParams.Email,
		HashedPasswords: newHashedPassword,
		ID: user.ID,
	}
	updatedUser, err := cfg.db.ChangePassEmail(req.Context(), dbParams)
	if err != nil {
		respondWithError(resWrite, http.StatusForbidden, "err changing user vals:", err)
		return
	}
	resParams := UserCreate{
		ID: updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email: updatedUser.Email,
	}
	respondWithJSON(resWrite, http.StatusOK, resParams)
}