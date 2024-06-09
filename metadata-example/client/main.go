package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/madeindra/grpc-metadata/metadata-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// buat koneksi ke server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// buat client
	client := pb.NewHelloServiceClient(conn)

	// panggil unary rpc
	callUnary(client)

	// panggil stream rpc
	callStream(client)
}

func callUnary(client pb.HelloServiceClient) {
	// buat metadata
	metadataFromClient := metadata.New(map[string]string{"x-client-id": "my-client"})

	// buat context yang sudah diisi metadata
	ctx := metadata.NewOutgoingContext(context.Background(), metadataFromClient)

	// siapkan variable untuk menerima metadata & trailer dari server
	var metadataFromServer, trailerFromServer metadata.MD

	// panggil greet rpc
	resp, err := client.Greet(ctx, nil, grpc.Header(&metadataFromServer), grpc.Trailer(&trailerFromServer))
	if err != nil {
		log.Fatalf("error when calling Greet: %v", err)
	}

	// tampilkan response dari server
	fmt.Println("---Response from Greet---")
	fmt.Println(resp.String())
	fmt.Println()

	// tampilkan metadata dari server
	fmt.Println("---Metadata from Greet---")
	for key, val := range deduplicateMD(metadataFromServer) {
		fmt.Println(key + ":" + val)
	}
	fmt.Println()

	// tampilkan trailer dari server
	fmt.Println("---Trailer from Greet---")
	for key, val := range deduplicateMD(trailerFromServer) {
		fmt.Println(key + ":" + val)
	}
	fmt.Println()
}

func callStream(client pb.HelloServiceClient) {
	// panggil stream rpc
	stream, err := client.SeverTime(context.Background(), nil)
	if err != nil {
		log.Fatalf("error when calling SeverTime: %v", err)
	}

	// baca stream
	fmt.Println("---Response from SeverTime---")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println()
			break
		} else if err != nil {
			log.Fatalf("error when reading stream: %v", err)
		}

		// tampilkan response dari server
		fmt.Println(resp.String())
	}

	// tampilkan metadata dari server
	fmt.Println("---Metadata from SeverTime---")
	metadataFromServer, err := stream.Header()
	if err != nil {
		log.Fatalf("error when reading metadata: %v", err)
	}

	for key, val := range deduplicateMD(metadataFromServer) {
		fmt.Println(key + ":" + val)
	}
	fmt.Println()

	// tampilkan trailer dari server
	fmt.Println("---Trailer from SeverTime---")
	for key, val := range deduplicateMD(stream.Trailer()) {
		fmt.Println(key + ":" + val)
	}
	fmt.Println()
}

func deduplicateMD(md metadata.MD) map[string]string {
	metadataReceived := make(map[string]string)
	for key, val := range md {
		metaValue := ""
		for _, v := range val {
			metaValue += v + " "
		}

		metadataReceived[key] = metaValue
	}

	return metadataReceived
}
