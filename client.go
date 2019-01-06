package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func RunClient() {
	address := ":12443"

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Unable to connect to GRPC server (", address, "):", err)
	} else {
		log.Println("Connected to GRPC server @ ", address)
	}

	tsc := NewTunnelServiceClient(conn)

	tc, err := tsc.Tunnel(context.Background())
	if err != nil {
		panic(err)
	}

	err = tc.Send(&ChunkOrSetup{X: &ChunkOrSetup_SetupRequest{&SetupRequest{
		Ip:   "127.0.0.1",
		Port: 11180,
	}}})

	if err != nil {
		panic(err)
	}

	sendPacket := func(data []byte) error {
		return tc.Send(&ChunkOrSetup{X: &ChunkOrSetup_Data{data}})
	}

	sendPacket([]byte("GET / HTTP/1.0\r\nHost: hiber.global\r\nUser-agent: Dank-User-Agent-Erwin\r\n"))

	for {
		chunk, err := tc.Recv()
		if err != nil {
			log.Println("Recv terminated:", err)
			return
		}

		println(string(chunk.Data))
	}
}
