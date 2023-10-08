package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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

func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Currently in the check content type middleware")

		//filtering request by MIME type
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 -- Unsupported Media Type. Please input JSON!"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		//Setting cookie to every API response
		cookie := http.Cookie{Name: "Server-Time(UTC)", Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("Currently in the set server time middleware")
	})
}
func main() {
	originalHandler := http.HandlerFunc(postHandler)

	http.Handle("/city", filterContentType(setServerTimeCookie(originalHandler)))
	http.ListenAndServe(":8000", nil)
}
