package main

import (
	"context"
	"time"

	pb "github.com/madeindra/grpc-metadata/metadata/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) Greet(ctx context.Context, _ *emptypb.Empty) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{
		Metadata: map[string]string{
			"Hello": "World",
		},
	}, nil
}

func (s *server) SeverTime(_ *emptypb.Empty, stream pb.HelloService_SeverTimeServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			currentTime := time.Now()

			if err := stream.Send(&pb.ServerTimeResponse{
				CurrentTime: timestamppb.New(currentTime),
			}); err != nil {
				return err
			}

			time.Sleep(1 * time.Second)
		}
	}
}
