package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type city struct {
	Name string
	Area uint64
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var reqCity city
		err := json.NewDecoder(r.Body).Decode(&reqCity)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Got %s city with area %d sq. miles.\n", reqCity.Name, reqCity.Area)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 --- Created Successfully"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 --- Method Not Allowed"))
	}

}
func main() {
	http.HandleFunc("/city", postHandler)
	http.ListenAndServe(":8000", nil)

}
