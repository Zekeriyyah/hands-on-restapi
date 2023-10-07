package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category is: %v\n", vars["category"])
	fmt.Fprintf(w, "ID is: %v\n", vars["id"])
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()

	router.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 1 * time.Millisecond,
		ReadTimeout:  1 * time.Millisecond,
	}

	fmt.Println("Server successfully running on port ", port)
	log.Fatal(srv.ListenAndServe())
}
