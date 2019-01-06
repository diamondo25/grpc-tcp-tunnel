package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
	"os"
)

type GrpcServer struct{}

var grpcServer = &GrpcServer{}

func (s *GrpcServer) Tunnel(ts TunnelService_TunnelServer) error {
	md, ok := metadata.FromIncomingContext(ts.Context())
	if !ok {
		return fmt.Errorf("unable to get metadata from context")
	}

	connectIpList := md.Get("connect_ip")
	if len(connectIpList) < 1 {
		return fmt.Errorf("expected connect_ip in metadata")
	}
	connectIp := connectIpList[0]

	connectPortList := md.Get("connect_port")
	if len(connectPortList) < 1 {
		return fmt.Errorf("expected connect_port in metadata")
	}
	connectPort := connectPortList[0]

	addr := net.JoinHostPort(connectIp, connectPort)

	log.Println("Connecting to", addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	// Make sure we close it
	defer conn.Close()
	defer log.Println("Connection closed to", addr)

	errChan := make(chan error)

	// Writing loop
	go func() {
		for {
			c, err := ts.Recv()
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error while receiving data:", err)
				}
				errChan <- nil
				return
			}

			data := c.Data

			fmt.Println("Sending bytes to tcp server", len(data))

			_, err = conn.Write(data)
			if err != nil {
				errChan <- fmt.Errorf("unable to write to connection: %v", err)
				return
			}
		}
	}()

	// Reading loop
	go func() {
		buff := make([]byte, 10000)
		for {
			bytesRead, err := conn.Read(buff)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error while receiving data:", err)
				} else {
					fmt.Println("Remote connection closed")
				}
				errChan <- nil
				return
			}

			fmt.Println("Sending bytes to grpc client", bytesRead)

			err = ts.Send(&Chunk{
				Data: buff[0:bytesRead],
			})

			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Blocking read
	returnedError := <-errChan

	return returnedError
}

func RunServer() {
	address := os.Args[2]
	x := grpc.NewServer()
	RegisterTunnelServiceServer(x, grpcServer)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Unable to setup socket on", address, ":", err)
	}

	log.Print("Starting GRPC server on ", address, "...")
	if err := x.Serve(lis); err != nil {
		log.Fatal("Unable to serve:", err)
	}
}
