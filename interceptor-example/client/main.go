package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/fastcampus-backend-golang/grpc-metadata-example/interceptor-example/proto"
)

func main() {
	// buat koneksi ke server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// buat client
	client := pb.NewSecretServiceClient(conn)

	// panggil rpc pembuat token
	token := createToken(client)

	// panggil unary rpc tanpa token
	callUnaryWithToken(client, "")

	// panggil unary rpc yang memerlukan token
	callUnaryWithToken(client, token)

	// panggil stream rpc tanpa token
	callStreamWithToken(client, "")

	// panggil stream rpc yang memerlukan token
	callStreamWithToken(client, token)
}

func createToken(client pb.SecretServiceClient) string {
	// panggil rpc pembuat token
	resp, err := client.Token(context.Background(), nil)
	if err != nil {
		log.Fatalf("error when calling CreateToken: %v", err)
	}

	// tampilkan response dari server
	fmt.Println("---Response from CreateToken---")
	fmt.Println(resp.String())
	fmt.Println()

	// kembalikan token
	return resp.Token
}

func callUnaryWithToken(client pb.SecretServiceClient, token string) {
	// buat metadata
	metadataFromClient := metadata.New(map[string]string{"token": token})

	// buat context yang sudah diisi metadata
	ctx := metadata.NewOutgoingContext(context.Background(), metadataFromClient)

	// panggil greet rpc
	resp, err := client.Protected(ctx, nil)
	if err != nil {
		fmt.Println("---Error when calling Protected---")
		fmt.Println("Token used:", token)
		fmt.Println(err)
		fmt.Println()

		return
	}

	// tampilkan response dari server
	fmt.Println("---Response from Protected---")
	fmt.Println("Token used:", token)
	fmt.Println(resp.String())
	fmt.Println()
}

func callStreamWithToken(client pb.SecretServiceClient, token string) {
	// buat metadata
	metadataFromClient := metadata.New(map[string]string{"token": token})

	// buat context yang sudah diisi metadata
	ctx := metadata.NewOutgoingContext(context.Background(), metadataFromClient)

	// panggil stream rpc
	stream, err := client.ProtectedStream(ctx, nil)
	if err != nil {
		// di titik ini, stream berhasil dibuat dan request akan masuk ke interceptor

		return
	}

	// baca stream
	fmt.Println("---Response or Error from ProtectedStream---")
	fmt.Println("Token used:", token)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println()
			break
		}
		if err != nil {
			fmt.Println(err)
			fmt.Println()

			break
		}

		fmt.Println(resp.String())
	}

	fmt.Println()
}
