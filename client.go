package main

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
)

func RunClient() {
	address := os.Args[2]

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Unable to connect to GRPC server (", address, "):", err)
	}

	md := metadata.Pairs(
		"connect_ip", os.Args[3],
		"connect_port", os.Args[4],
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tsc := NewTunnelServiceClient(conn)

	tc, err := tsc.Tunnel(ctx)
	if err != nil {
		panic(err)
	}

	sendPacket := func(data []byte) error {
		return tc.Send(&Chunk{Data: data})
	}

	go func() {
		for {
			chunk, err := tc.Recv()
			if err != nil {
				log.Println("Recv terminated:", err)
				os.Exit(0)
			}

			os.Stdout.Write(chunk.Data)
		}
	}()

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		sendPacket(buf)
	}
}
