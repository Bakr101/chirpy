package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(resWrite http.ResponseWriter, req *http.Request){
	cfg.fileserverHits.Store(0)
	resWrite.WriteHeader(http.StatusOK)
	resWrite.Write([]byte(fmt.Sprintf("Metrics Reset. Hits: %v", cfg.fileserverHits.Load())))
	
}