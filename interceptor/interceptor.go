package main

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var allowedToken = map[string]time.Time{}

const tokenDuration = 30 * time.Second

func validateToken(token string) bool {
	expiredAt, ok := allowedToken[token]
	if !ok {
		return false
	}

	if time.Now().After(expiredAt) {
		return false
	}

	return true
}

func ProtectedInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// baca rpc yang dipanggil
	method := info.FullMethod

	// whitelist rpc dari interceptor
	if method == "/metadata.SecretService/Token" {
		return handler(ctx, req)
	}

	// ambil metada
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("token required")
	}

	// cek jumlah metadata token
	if len(meta["token"]) < 1 {
		return nil, errors.New("token invalid")
	}

	// gunakan hanya token pertama
	if !validateToken(meta["token"][0]) {
		return nil, errors.New("expired or invalid token")
	}

	return handler(ctx, req)
}
