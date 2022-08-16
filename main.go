package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"tcp_service/servers"
)

func main() {

	addr := *flag.String("address", "localhost:8080", "server address")
	saddr := *flag.String("sender_address", ":5005", "drpc sender address")
	flag.Parse()
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("TCP Server has started")

	for {
		connect, err := server.Accept()

		if err != nil {
			log.Fatal(err)
			return
		}
		go servers.Save(connect, saddr)
	}
}
