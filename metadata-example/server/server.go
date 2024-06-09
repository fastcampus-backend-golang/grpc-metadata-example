package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/fastcampus-backend-golang/grpc-metadata-example/metadata-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) Greet(ctx context.Context, _ *emptypb.Empty) (*pb.GreetResponse, error) {
	// tentukan waktu sekarang
	currentTime := time.Now()

	// defer pengiriman trailer dari server (agar dikirim setelah response)
	defer func() {
		serverTrailer := metadata.New(map[string]string{"x-process-time": fmt.Sprintf("%dms", time.Since(currentTime).Microseconds())})
		grpc.SetTrailer(ctx, serverTrailer)
	}()

	// map response
	response := make(map[string]string)

	// baca metadata dari request
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for key, val := range meta {
			var value string

			for _, v := range val {
				// gabungkan value dari metadata yang sama
				value = fmt.Sprintf("%s%s", value, v)
			}

			response[key] = value
		}
	}

	// kirim metadata dari server (harus sebelum response)
	serverMetadata := metadata.New(map[string]string{"x-server-rpc": "greet-rpc"})
	grpc.SendHeader(ctx, serverMetadata)

	// kirim response
	return &pb.GreetResponse{
		Metadata: response,
	}, nil
}

func (s *server) SeverTime(_ *emptypb.Empty, stream pb.HelloService_SeverTimeServer) error {
	// tentukan waktu sekarang
	currentTime := time.Now()

	// defer pengiriman trailer
	defer func() {
		serverTrailer := metadata.New(map[string]string{"x-process-time": fmt.Sprintf("%dms", time.Since(currentTime).Microseconds())})
		stream.SetTrailer(serverTrailer)
	}()

	// kirim metadata dari server
	serverMetadata := metadata.New(map[string]string{"x-server-rpc": "server-time-rpc"})
	stream.SendHeader(serverMetadata)

	// proses pengiriman waktu
	for time.Since(currentTime) < 3*time.Second {
		select {
		case <-stream.Context().Done():
			// ketika dicancel oleh client, trailer tidak akan dikirim

			return stream.Context().Err()

		case <-time.After(1 * time.Second):
			// dilakukan setiap 1 detik

			if err := stream.Send(&pb.ServerTimeResponse{
				CurrentTime: timestamppb.New(time.Now()),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
