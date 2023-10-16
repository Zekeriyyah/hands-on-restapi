package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
)

func pingTime(req *restful.Request, resp *restful.Response) {
	//Write current time to the response
	io.WriteString(resp, fmt.Sprintf("%v", time.Now()))
}
func main() {
	//Create a web service
	webservice := new(restful.WebService)

	//Create a route and attach it to handler in the service
	webservice.Route(webservice.GET("/ping").To(pingTime))

	//Add the service to application
	restful.Add(webservice)
	http.ListenAndServe(":8000", nil)
}
