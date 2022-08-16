package servers

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"time"

	"google.golang.org/grpc"
)

type Capture struct {
	TimeStamp      time.Time     `json: "time"`
	CaptureLength  int           `json: "caplength"`
	Length         int           `json: "length"`
	InterfaceIndex int           `json :  "index"`
	AccalaryData   []interface{} `json: "accalary"`
}

type Packet struct {
	Ci   Capture
	Data []byte
}

func Save(connect net.Conn, addr string) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := NewClient(conn)

	packets := []Packet{}
	for {
		read, err := ReceiveALL(connect, 8)

		if err != nil {
			log.Fatal(err)
			return
		}

		size := binary.BigEndian.Uint64(read)
		fmt.Println(size)

		read, err = ReceiveALL(connect, size)

		if err != nil {
			log.Fatal(err)
			return
		}

		if size == 4 && string(read) == "STOP" {
			break
		}

		cap := Capture{}

		err = json.Unmarshal(read, &cap)
		read, err = ReceiveALL(connect, 8)

		if err != nil {
			log.Fatal(err)
			return
		}

		size = binary.BigEndian.Uint64(read)

		read, err = ReceiveALL(connect, size)

		if err != nil {
			log.Fatal(err)
			return
		}

		packets = append(packets, Packet{cap, read})

		fmt.Printf("File size: %v\n", size)

	}
	fmt.Println("Stopped receiving")

	statistics, err := client.Upload(context.Background(), packets)
	if err != nil {
		connect.Write([]byte("Could not make statistics"))
		connect.Close()
		fmt.Println("File receiving has ended")
		return
	}
	connect.Write([]byte(statistics))
	connect.Close()
	fmt.Println("File receiving has ended")
	fmt.Println()

}

func ReceiveALL(connect net.Conn, size uint64) ([]byte, error) {
	read := make([]byte, size)

	_, err := io.ReadFull(connect, read)
	if err != nil {
		log.Fatal(err)
		return []byte{}, err
	}

	return read, nil
}
