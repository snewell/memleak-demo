package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	"google.golang.org/grpc/xds"

	"github.com/spf13/cobra"

	"github.com/snewell/memleak-demo/internal/pb"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start an xds-enabled grpc server",
		Run:   serverFn,
	}

	listenPort int
)

type ss struct {
	pb.UnimplementedSSServer
}

func (s ss) DoIt(*pb.Request, grpc.ServerStreamingServer[pb.Response]) error {
	return nil
}

func serverFn(cmd *cobra.Command, args []string) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	listenAddr := fmt.Sprintf(":%v", listenPort)
	log.Printf("Listening on %v", listenAddr)
	lis, err := net.Listen("tcp", listenAddr)
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
	log.Printf("Starting server...")
	if err := xdsServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	serverCmd.Flags().IntVar(&listenPort, "port", 50051, "Listening port")

	rootCmd.AddCommand(serverCmd)
}
