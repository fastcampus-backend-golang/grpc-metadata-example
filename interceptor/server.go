package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	pb "github.com/madeindra/grpc-metadata/interceptor/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var allowedNonce = map[string]time.Time{}

const nonceDuration = 10 * time.Second

type server struct {
	pb.UnimplementedSecretServiceServer
}

func (*server) Nonce(ctx context.Context, _ *emptypb.Empty) (*pb.NonceResponse, error) {
	// buat kunci baru
	genValue := fmt.Sprintf("%d", rand.Int())
	allowedNonce[genValue] = time.Now().Add(nonceDuration)

	// kirim kunci dalam response
	return &pb.NonceResponse{
		Nonce: genValue,
	}, nil
}

func (*server) Protected(ctx context.Context, _ *emptypb.Empty) (*pb.ProtectedResponse, error) {
	return &pb.ProtectedResponse{
		Message: "Welcome!",
	}, nil
}
