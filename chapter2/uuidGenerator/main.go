package main

import (
	"log"
	"net/http"

	multiplexhandlers "github.com/Zekeriyyah/hands-on-RestApis/chapter2/uuidGenerator/multiplexHandlers"
)

func main() {

	mux := &multiplexhandlers.UUID{}
	log.Fatal(http.ListenAndServe(":8080", mux))
}
