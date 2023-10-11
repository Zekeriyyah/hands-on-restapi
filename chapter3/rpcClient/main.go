package main

import (
	"log"
	"net/rpc"
)

type Args struct{}

func main() {
	var reply int64
	args := Args{}

	//Client dialing an RPC server
	client, err := rpc.DialHTTP("tcp", "localhost"+":1234")
	if err != nil {
		log.Fatal("Dialing-Error--", err)
		return
	}

	//Client making a remote call to GiveServerTime method of TimeServer
	err = client.Call("TimeServer.GiveServerTime", args, &reply)
	if err != nil {
		log.Fatal("CallingError-- ", err)
		return
	}

	log.Printf("%d", reply)
}
