package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Args struct{}

type TimeServer int64

func (t *TimeServer) GiveServerTime(args *Args, reply *int64) error {
	//Fill reply pointer with current time to send the data back
	*reply = time.Now().Unix()
	return nil
}

func main() {
	timeServer := new(TimeServer)
	rpc.Register(timeServer)

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Error-Listening: ", e)
		return
	}
	fmt.Println("Server is running....")
	http.Serve(l, nil)
}
