package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/Bakr101/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserCreate struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool	`json:"is_chirpy_red"`
	
}

type UserReq struct{
	Email string 
	Password string 
}
func (cfg *apiConfig)handlerCreateUser(resWrite http.ResponseWriter, req *http.Request){
	
	reqParams := UserReq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
	}
	hashed_pass,err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	dbParams := database.CreateUserParams{
		Email: reqParams.Email,
		HashedPasswords: hashed_pass,
	}
	user, err := cfg.db.CreateUser(context.Background(), dbParams)
	if err != nil {
		//log.Fatalf("error creating user err:%s", err)
		respondWithError(resWrite, http.StatusInternalServerError, "error creating user", err)
		return
	}
	userJson := UserCreate{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(resWrite, http.StatusCreated, userJson)
}