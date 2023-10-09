package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing Request....")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	log.Println("Finished processing request....")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", handle)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":8000", loggedRouter)
	/*** Other useful middlewares in gorilla/handlers includes
	1. CompressionHandler- for zipping the responses
	2. RecoveryHandler- for recovering from unexpected panics
	**/

}
