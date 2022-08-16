package servers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	uploadpb "tcp_service/proto"

	"google.golang.org/grpc"
)

type Client struct {
	client uploadpb.UploadServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: uploadpb.NewUploadServiceClient(conn),
	}
}

func (c Client) Upload(con context.Context, packets []Packet) (string, error) {
	ctx, cancel := context.WithDeadline(con, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := c.client.Upload(ctx)
	if err != nil {

		return "", err
	}

	for _, pack := range packets {

		b, err := json.Marshal(&pack.Ci)

		if err != nil {

			return "", err
		}

		if err := stream.Send(&uploadpb.UploadRequest{Chunk: b}); err != nil {

			return "", err
		}

		if err := stream.Send(&uploadpb.UploadRequest{Chunk: pack.Data}); err != nil {

			return "", err
		}

	}

	res, err := stream.CloseAndRecv()
	if err != nil {

		return "", err
	}
	fmt.Println("stopped sending")

	return res.GetName(), nil
}
