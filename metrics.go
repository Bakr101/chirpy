package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerHits(resWrite http.ResponseWriter, req *http.Request){
	html := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())

	resWrite.Header().Add("Content-Type", "text/html")
	resWrite.WriteHeader(http.StatusOK)
	resWrite.Write([]byte(html))
	
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       cfg.fileserverHits.Add(1)
	   next.ServeHTTP(w, r)
    })
}