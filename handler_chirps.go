package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

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

func (cfg *apiConfig)handlerCreateChirp(resWrite http.ResponseWriter, req *http.Request){
	type chirpReq struct {
		Body string `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}
	
	decoder := json.NewDecoder(req.Body)
	reqParams := chirpReq{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	const maxChirpLength = 140
	if len(reqParams.Body) > maxChirpLength {
		respondWithError(resWrite, http.StatusBadRequest,"Chirp is too long", nil)
		return	
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: reqParams.Body,
		UserID: reqParams.User_ID,
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