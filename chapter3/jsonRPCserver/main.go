package main

import (
	jsonparse "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Args struct {
	ID string
}
type Book struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Author string `json:"author,omitempty"`
}

type JSONServer struct{}

func (j *JSONServer) GiveBookDetail(r *http.Request, args *Args, reply *Book) error {
	absPath, _ := filepath.Abs("books.json")

	//Reading JSON data
	raw, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatal("Read-JSON-ERROR: ", err)
		return err
	}

	//Unmarshaling the raw JSON data
	var books []Book
	marshalErr := jsonparse.Unmarshal(raw, &books)
	if marshalErr != nil {
		log.Fatal("ERROR-JSONUnmarshal: ", marshalErr)
		return marshalErr
	}

	//Finding the requested book and filling the details into reply pointer
	for _, book := range books {
		if book.ID == args.ID {
			*reply = book
			break
		}
	}
	return nil
}
func main() {
	s := rpc.NewServer()

	//Registering codec type to serve
	s.RegisterCodec(json.NewCodec(), "application/json")

	//Registering service
	s.RegisterService(new(JSONServer), "")

	r := mux.NewRouter()
	r.Handle("/rpc", s)

	fmt.Println("Server running...")
	log.Fatal(http.ListenAndServe(":1234", r))

}
