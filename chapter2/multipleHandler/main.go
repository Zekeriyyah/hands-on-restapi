package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {

	newMux := http.NewServeMux()

	newMux.HandleFunc("/random-float",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, rand.Float64())
		})
	newMux.HandleFunc("/random-int",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, rand.Intn(100))
		})

	log.Fatal(http.ListenAndServe(":8080", newMux))
}
