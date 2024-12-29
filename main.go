package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
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