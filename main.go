package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	//constants
	const port = "8080"
	const filepathRoot = "."
	//Serve Mux (Router)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", fs))
	mux.HandleFunc("/healthz", handlerReadiness)
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