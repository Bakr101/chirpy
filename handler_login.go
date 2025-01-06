package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/Bakr101/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserLogin struct{
	ID        		uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	Email     		string    	`json:"email"`
	Token     		string		`json:"token"`
	RefreshToken 	string 		`json:"refresh_token"`
}

//Time to refresh In days
const refreshTokenTime = 60

func (cfg *apiConfig)handlerLogin(resWrite http.ResponseWriter, req *http.Request){
	type userReq struct{
		Email string 
		Password string 
		
	}
	reqParams := userReq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
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
	JWTToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(1))
	if err != nil{
		respondWithError(resWrite, http.StatusInternalServerError,"error creating token", err)
	}

	//Create Referesh Token & save to DB
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error creating refresh token", err)
	}
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, refreshTokenTime),
	}

	refreshTokenDB, err := cfg.db.CreateRefreshToken(req.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error inserting into refresh_tokens DB", err)
	}

	userRes := UserLogin{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: JWTToken,
		RefreshToken: refreshTokenDB.Token,
	}
	
	respondWithJSON(resWrite, http.StatusOK, userRes)
}