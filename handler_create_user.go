package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig)handlerCreateUser(resWrite http.ResponseWriter, req *http.Request){
	type userReq struct{
		Email string `json:"email"`
	}
	
	reqParams := userReq{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error decoding json", err)
		return
	}
	user, err := cfg.db.CreateUser(context.Background(), reqParams.Email)
	if err != nil {
		//log.Fatalf("error creating user err:%s", err)
		respondWithError(resWrite, http.StatusInternalServerError, "error creating user", err)
		return
	}
	userJson := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(resWrite, http.StatusCreated, userJson)
}