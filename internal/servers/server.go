package servers

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"tcp_service/internal/metrics"

	"google.golang.org/grpc"
)

func Save(connect net.Conn, addr string) {
	metrics.RecordMetrics()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := NewClient(conn)
	fileConntent := []byte{}

	for {
		read, err := ReceiveALL(connect, 8)

		if err != nil {
			log.Fatal(err)
			return
		}

		size := binary.BigEndian.Uint64(read)

		read, err = ReceiveALL(connect, size)

		if err != nil {
			log.Fatal(err)
			return
		}

		if size == 4 && string(read) == "STOP" {

			break
		}

		fileConntent = append(fileConntent, read...)

		fmt.Printf("File size: %v\n", size)

	}
	fmt.Println("Stopped receiving")

	statistics, err := client.Upload(context.Background(), fileConntent)
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
