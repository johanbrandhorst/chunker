package main

import (
	"context"
	"io"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/johanbrandhorst/chunker/protos/chunker"
)

func main() {
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cc := chunker.NewChunkerClient(conn)
	client, err := cc.Chunker(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	var blob []byte
	for {
		c, err := client.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("Transfer of %d bytes successful", len(blob))
				return
			}

			panic(err)
		}

		blob = append(blob, c.Chunk...)
	}
}
