package main

import "net/http"

func handlerReadiness(resWrite http.ResponseWriter, req *http.Request){
	resWrite.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWrite.WriteHeader(http.StatusOK)
	resWrite.Write([]byte(http.StatusText(http.StatusOK)))
}