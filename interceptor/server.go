package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	pb "github.com/madeindra/grpc-metadata/interceptor/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedSecretServiceServer
}

func (*server) Token(ctx context.Context, _ *emptypb.Empty) (*pb.TokenResponse, error) {
	// buat kunci baru
	genValue := fmt.Sprintf("%d", rand.Int())
	allowedToken[genValue] = time.Now().Add(tokenDuration)

	// kirim kunci dalam response
	return &pb.TokenResponse{
		Token: genValue,
	}, nil
}

func (*server) Protected(ctx context.Context, _ *emptypb.Empty) (*pb.ProtectedResponse, error) {
	return &pb.ProtectedResponse{
		Message: "Welcome!",
	}, nil
}
