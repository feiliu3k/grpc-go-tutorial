// Package main implements a client for Echo service.
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	pb "github.com/wangy8961/grpc-go-tutorial/features/echopb"
	"google.golang.org/grpc"
)

type basicAuth struct {
	username string
	password string
}

// Return value is mapped to request headers.
func (b *basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := b.username + ":" + b.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

// 是否使用 TLS 安全加密
func (b *basicAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	certFile := flag.String("cacert", "cacert.pem", "CA root certificate")
	username := flag.String("username", "admin", "The username to authenticate with")
	password := flag.String("password", "password", "The password to authenticate with")
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile(*certFile, "")
	if err != nil {
		log.Fatalf("failed to load CA root certificate: %v", err)
	}

	opts := []grpc.DialOption{
		// 1. TLS 认证
		grpc.WithTransportCredentials(creds),
		// 2. basic token 认证
		grpc.WithPerRPCCredentials(&basicAuth{
			username: *username,
			password: *password,
		}),
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, opts...) // To call service methods, we first need to create a gRPC channel to communicate with the server. We create this by passing the server address and port number to grpc.Dial()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn) // Once the gRPC channel is setup, we need a client stub to perform RPCs. We get this using the NewEchoClient method provided in the pb package we generated from our .proto.

	// Contact the server and print out its response.
	msg := "madmalls.com"
	resp, err := c.UnaryEcho(context.Background(), &pb.EchoRequest{Message: msg}) // Now let’s look at how we call our service methods. Note that in gRPC-Go, RPCs operate in a blocking/synchronous mode, which means that the RPC call waits for the server to respond, and will either return a response or an error.
	if err != nil {
		log.Fatalf("failed to call UnaryEcho: %v", err)
	}
	fmt.Printf("response:\n")
	fmt.Printf(" - %q\n", resp.GetMessage())
}
