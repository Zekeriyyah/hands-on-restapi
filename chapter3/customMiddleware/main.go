package main

import (
	"fmt"
	"net/http"
)

func middleware(origHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middleware logic before request phase")
		origHandler.ServeHTTP(w, r)
		fmt.Println("Executing middleware logic after response phase")
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing main handler")
	w.Write([]byte("OK"))
}

func main() {
	originalHandler := http.HandlerFunc(handle)

	http.Handle("/", middleware(originalHandler))

	http.ListenAndServe(":8000", nil)
}
