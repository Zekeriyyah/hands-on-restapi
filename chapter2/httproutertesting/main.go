package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/api/v1/version", goVersion)
	router.GET("/api/v1/show-file/:name", getFileContent)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getCommandOutput(command string, arguments ...string) string {
	out, _ := exec.Command(command, arguments...).Output()
	return string(out)
}

func goVersion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	response := getCommandOutput("/usr/local/go/bin/go", "version")
	io.WriteString(w, response)
}

func getFileContent(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, getCommandOutput("/usr/bin/cat", params.ByName("name")))
}
