# gRPC Metadata Example

1. Clone repository ini
```bash
git clone git@github.com:fastcampus-backend-golang/grpc-metadata-example.git
```

2. Masuk ke direktori
```bash
cd grpc-metadata-example
```

Jalankan server Metadata dan Interceptor dalam waktu yang terpisah, karena kedua server berjalan di port yang sama, yaitu `:50051`.

Menjalankan keduanya di saat yang bersamaan akan menyebabkan konflik port.

## Cara Menjalankan Contoh Metadata

1. Jalankan server metadata dengan perintah
```bash
make run-server-metadata
```

3. Jalankan client metadata dengan perintah
```bash
make run-client-metadata
```

## Cara Menjalankan Contoh Interceptor
1. Jalankan server interceptor dengan perintah
```bash
make run-server-interceptor
```

2. Jalankan client interceptor dengan perintah
```bash
make run-client-interceptor
```

## Konten
- metadata-example: Contoh penggunaan metadata (header) pada gRPC server & client
- interceptor-example: Contoh penggunaan interceptor metadata (middleware) pada gRPC server & client
 