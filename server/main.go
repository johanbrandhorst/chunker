package main

import (
	"log"
	"math/rand"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/johanbrandhorst/chunker/protos/chunker"
)

const chunkSize = 64 * 1024 // 64 KiB

type chunkerSrv []byte

func (c chunkerSrv) Chunker(_ *empty.Empty, srv chunker.Chunker_ChunkerServer) error {
	chnk := &chunker.Chunk{}
	for currentByte := 0; currentByte < len(c); currentByte += chunkSize {
		if currentByte+chunkSize > len(c) {
			chnk.Chunk = c[currentByte:len(c)]
		} else {
			chnk.Chunk = c[currentByte : currentByte+chunkSize]
		}
		if err := srv.Send(chnk); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}

	g := grpc.NewServer()
	blob := make([]byte, 128*1024*1024) // 128MiB
	rand.Read(blob)
	chunker.RegisterChunkerServer(g, chunkerSrv(blob))

	log.Println("Serving on localhost:10000")
	log.Fatalln(g.Serve(lis))
}
