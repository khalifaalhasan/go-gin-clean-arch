# --- Stage 1: Build (Kita masak aplikasinya) ---
FROM golang:1.23-alpine AS builder

# Set folder kerja di dalam container
WORKDIR /app

# Copy file dependency dulu (biar cache-nya optimal)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code Go
COPY . .

# Build aplikasi jadi binary namanya 'main'
# cmd/api/main.go adalah lokasi file main.go Anda
RUN go build -o main cmd/api/main.go

# --- Stage 2: Run (Kita sajikan di piring bersih) ---
FROM alpine:latest

WORKDIR /root/

# Copy hasil build dari Stage 1
COPY --from=builder /app/main .

# Expose port (hanya dokumentasi, Render nanti override via ENV)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]