package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	fileserverHits	atomic.Int32
}

func main(){
	//constants
	const port = "8080"
	const filepathRoot = "."
	
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	//Serve Mux (Router)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fs)))

	//handling functions
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	//Server config
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

