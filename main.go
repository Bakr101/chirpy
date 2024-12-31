package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

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
	mux.HandleFunc("GET /admin/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	
	
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}

func handlerReadiness(resWrite http.ResponseWriter, req *http.Request){
	resWrite.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWrite.WriteHeader(200)
	resWrite.Write([]byte("OK"))
}

type apiConfig struct{
	fileserverHits	atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       cfg.fileserverHits.Add(1)
	   next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) handlerHits(resWrite http.ResponseWriter, req *http.Request){
	html := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())
	resWrite.Header().Set("Content-Type", "text/html")
	resWrite.WriteHeader(200)
	resWrite.Write([]byte(html))
	
}

func (cfg *apiConfig) handlerReset(resWrite http.ResponseWriter, req *http.Request){
	cfg.fileserverHits.Swap(0)
	resWrite.Write([]byte(fmt.Sprintf("Metrics Reset. Hits: %v", cfg.fileserverHits.Load())))
	
}

func handlerValidateChirp(resWrite http.ResponseWriter, req *http.Request){
	type chirpReq struct{
		Body string `json:"body"`
	}
	type chirpResSuccess struct {
		Valid bool `json:"valid"`
	}
	type chirpResErr struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := chirpReq{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		log.Printf("Error decoding request params: %s", err)
		resWrite.WriteHeader(500)
		return
	}
	
	if len(reqParams.Body) > 140 {
		
		respBody := chirpResErr{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error encoding response body: %s", err)
			resWrite.WriteHeader(500)
			return
		}
		resWrite.WriteHeader(400)
		resWrite.Write(dat)
	}else{
		
		respBody := chirpResSuccess{
			Valid: true,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error encoding response body: %s", err)
			resWrite.WriteHeader(500)
			return
		}
		
		resWrite.WriteHeader(200)
		resWrite.Write(dat)
		return
	}
	
}