package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type GrpcServer struct{}

var grpcServer = &GrpcServer{}

func (s *GrpcServer) Tunnel(ts TunnelService_TunnelServer) error {
	initialReq, err := ts.Recv()
	if err != nil {
		return err
	}
	sr := initialReq.GetSetupRequest()

	if sr == nil {
		return fmt.Errorf("expected a SetupRequest")
	}

	addr := fmt.Sprintf("%s:%d", sr.Ip, sr.Port)

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

			data := c.GetData()
			if data == nil {
				errChan <- fmt.Errorf("expected data buffer")
				return
			}

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
	address := ":12443"
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
