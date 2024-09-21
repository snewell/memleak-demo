package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscreds "google.golang.org/grpc/credentials/xds"
	_ "google.golang.org/grpc/xds"

	"github.com/snewell/memleak-demo/internal/pb"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Perform busy work on a memleak-demo server",
		Run:   clientFn,
	}

	remoteHost string
	remotePort int

	workers      int
	requestCount int
)

func doWork(remote string) {
	creds, err := xdscreds.NewClientCredentials(xdscreds.ClientOptions{
		FallbackCreds: insecure.NewCredentials(),
	})
	if err != nil {
		log.Fatalf("error creating credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial(remote, opts...)
	if err != nil {
		log.Fatalf("error dialing: %v", err)
	}

	stub := pb.NewSSClient(conn)
	for i := 0; i < requestCount; i++ {
		stream, err := stub.DoIt(context.Background(), &pb.Request{})
		if err != nil {
			log.Printf("error doing scan: %v", err)
		} else {
			for {
				_, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("response error: %v", err)
				}
			}
		}
	}
}

func clientFn(cmd *cobra.Command, args []string) {
	remote := fmt.Sprintf("xds:///%v:%v", remoteHost, remotePort)
	log.Printf("Connecting to %v", remote)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			doWork(remote)
		}()
	}
	wg.Wait()
}

func init() {
	clientCmd.Flags().StringVar(&remoteHost, "host", "memleak-demo", "Remote host to connect to")
	clientCmd.Flags().IntVar(&remotePort, "port", 50051, "Remote port to connect to")
	clientCmd.Flags().IntVar(&workers, "workers", 10, "Number of workers (connections)")
	clientCmd.Flags().IntVar(&requestCount, "requests", 20, "Number of requests each worker should make")

	rootCmd.AddCommand(clientCmd)
}
