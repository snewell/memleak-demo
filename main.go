package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	"google.golang.org/grpc/xds"

	"github.com/snewell/memleak-demo/internal/pb"
)

type ss struct {
	pb.UnimplementedSSServer
}

func (s ss) DoIt(*pb.Request, grpc.ServerStreamingServer[pb.Response]) error {
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50551")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds := insecure.NewCredentials()
	if creds, err = xdscreds.NewServerCredentials(xdscreds.ServerOptions{FallbackCreds: creds}); err != nil {
		log.Fatalf("failed to create server-side xDS credentials: %v", err)
	}

	xdsServer, err := xds.NewGRPCServer(grpc.Creds(creds))
	if err != nil {
		log.Fatalf("failed to create xds server: %v", err)
	}
	s := ss{}
	pb.RegisterSSServer(xdsServer, s)
	if err := xdsServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
