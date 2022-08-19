package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	filename := *flag.String("fileName", "lo.pcapng", "pcap file directory")
	addr := *flag.String("address", "localhost:8080", "address of the GRPC server")
	flag.Parse()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	connect, err := net.Dial("tcp", addr)

	for {

		bin := make([]byte, 1024)
		n, err := file.Read(bin)

		if err != nil {
			log.Fatal(err)
			break
		}

		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(len(bin[:n])))

		_, err = connect.Write(b)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = connect.Write(bin[:n])
		if err != nil {
			log.Fatal(err)
			return
		}
		if n < 1024 {

			bin := make([]byte, 8)
			binary.BigEndian.PutUint64(bin, 4)
			connect.Write(bin)
			connect.Write([]byte("STOP"))
			break
		}

	}

	read := make([]byte, 1024)

	_, err = connect.Read(read)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(read))

	connect.Close()

}
