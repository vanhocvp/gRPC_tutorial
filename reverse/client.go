package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func doClientWork(clientChan <-chan *rpc.Client) {
	client := <-clientChan

	defer client.Close()

	var reply string

	err := client.Call("HelloService.Hello", "hello", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}

func main() {
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	clientChan := make(chan *rpc.Client)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Accept error:", err)
			}
			clientChan <- rpc.NewClient(conn)
		}
	}()

	doClientWork(clientChan)
}
