package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(resWrite http.ResponseWriter, req *http.Request){
	platform := cfg.platform
	if platform != "dev"{
		resWrite.WriteHeader(http.StatusForbidden)
		resWrite.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(resWrite, http.StatusInternalServerError, "couldn't reset users", err)
		return
	}
	cfg.fileserverHits.Store(0)
	resWrite.WriteHeader(http.StatusOK)
	resWrite.Write([]byte(fmt.Sprintf("Metrics Reset. Hits: %v  and database reset to initial state.", cfg.fileserverHits.Load())))
	
}