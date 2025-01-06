package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Bakr101/chirpy/internal/auth"
	"github.com/Bakr101/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct{
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_ID	  uuid.UUID	`json:"user_id"`
}

type Chirps struct{
	Chirps []Chirp `json:"chirps"`
}

func (cfg *apiConfig)handlerCreateChirp(resWrite http.ResponseWriter, req *http.Request){
	type chirpReq struct {
		Body string `json:"body"`
		
	}
	const maxChirpLength = 140
	
	decoder := json.NewDecoder(req.Body)
	reqParams := chirpReq{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	userToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "No Token Invalid Access", err)
		return
	}
	UserUUID, err := auth.ValidateJWT(userToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(resWrite, http.StatusUnauthorized, "Expired Token Invalid Access, please login to create a chirp", err)
		return
	}
	
	if len(reqParams.Body) > maxChirpLength {
		respondWithError(resWrite, http.StatusBadRequest,"Chirp is too long", nil)
		return	
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: reqParams.Body,
		UserID: UserUUID,
	})

	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "error creating chirp", err)
		return
	}
	chirpJson := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: handleBadwords(chirp.Body),
		User_ID: chirp.UserID,
	}
	respondWithJSON(resWrite, http.StatusCreated, chirpJson)

}

func handleBadwords(str string) string{
	splitted := strings.Split(str, " ")
	lowered := strings.ToLower(str)
	splittedToClean := strings.Split(lowered, " ")
	cleaned := []string{}
	badWords := map[string]string{
		"kerfuffle": "kerfuffle",
		"sharbert": "sharbert",
		"fornax": "fornax",
	}

	for idx, word := range splittedToClean{
		_, exists := badWords[word]
		if exists {
			cleaned = append(cleaned, "****")
		}
		if !exists {
			cleaned = append(cleaned, splitted[idx])
		}
	}

	return strings.Join(cleaned, " ")
}

