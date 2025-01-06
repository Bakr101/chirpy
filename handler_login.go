package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/google/uuid"
)

type UserLogin struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string	`json:"token"`
}

func (cfg *apiConfig)handlerLogin(resWrite http.ResponseWriter, req *http.Request){
	type userReq struct{
		Email string 
		Password string 
		ExpiresInSeconds int
	}
	reqParams := userReq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
	}
	//JWT Expiring Time
	if reqParams.ExpiresInSeconds > 3600 || reqParams.ExpiresInSeconds == 0{
		reqParams.ExpiresInSeconds = 3600
	}
	//Get User from DB & validate password
	user, err := cfg.db.GetUser(req.Context(), reqParams.Email)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "user email is invalid", err)
		return
	}
	err = auth.CheckPasswordHash(reqParams.Password, user.HashedPasswords)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "password is invalid", err)
		return
	}
	//Create a Token
	JWTToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(reqParams.ExpiresInSeconds))
	if err != nil{
		respondWithError(resWrite, http.StatusInternalServerError,"error creating token", err)
	}

	userRes := UserLogin{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: JWTToken,
	}
	
	respondWithJSON(resWrite, http.StatusOK, userRes)
}