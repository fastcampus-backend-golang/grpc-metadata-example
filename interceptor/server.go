package main

import (
	"context"

	pb "github.com/madeindra/grpc-metadata/interceptor/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedSecretServiceServer
}

func (*server) Nonce(ctx context.Context, _ *emptypb.Empty) (*pb.NonceResponse, error) {
	return &pb.NonceResponse{}, nil
}

func (*server) Protected(ctx context.Context, _ *emptypb.Empty) (*pb.ProtectedResponse, error) {
	return &pb.ProtectedResponse{}, nil
}
