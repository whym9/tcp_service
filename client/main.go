package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"encoding/json"
	"time"

	"github.com/google/gopacket/pcap"
)

type Capture struct {
	TimeStamp      time.Time     `json: "time"`
	CaptureLength  int           `json: "caplength"`
	Length         int           `json: "length"`
	InterfaceIndex int           `json :  "index"`
	AccalaryData   []interface{} `json: "accalary"`
}

func main() {
	filename := *flag.String("fileName", "lo.pcapng", "pcap file directory")
	addr := *flag.String("address", "localhost:8080", "address of the GRPC server")
	flag.Parse()
	handle, err := pcap.OpenOffline(filename)
	defer handle.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	connect, err := net.Dial("tcp", addr)

	for {

		data, ci, err := handle.ZeroCopyReadPacketData()
		if err == io.EOF || err != nil {
			bin := make([]byte, 8)
			binary.BigEndian.PutUint64(bin, 4)
			connect.Write(bin)
			connect.Write([]byte("STOP"))
			break
		}

		capi := Capture{
			ci.Timestamp,
			ci.CaptureLength,
			ci.Length,
			ci.InterfaceIndex,
			ci.AncillaryData,
		}

		b, _ := json.Marshal(&capi)

		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(len(b)))

		connect.Write(bin)
		connect.Write(b)

		bin = make([]byte, 8)

		binary.BigEndian.PutUint64(bin, uint64(len(data)))
		connect.Write(bin)
		connect.Write(data)

		if err != nil {
			panic(err)

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
