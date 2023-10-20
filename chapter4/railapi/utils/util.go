package utils

import (
	"log"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

func ParseID(response *restful.Response, id string) error {
	//Validating the Id in url to avoid sql injection

	_, parseErr := strconv.ParseInt(id, 0, 0)
	if parseErr != nil {
		log.Println("Invalid ID!: ", parseErr)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Invalid ID")
		return parseErr
	}
	return nil
}
