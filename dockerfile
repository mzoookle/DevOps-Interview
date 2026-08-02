# Stage 1: Build binary using Go
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Init modules and build the executable
COPY main.go .
RUN go mod init math-server && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Stage 2: Create Alpine image
FROM alpine:latest
WORKDIR /root/

# Copy binary from the builder stage
COPY --from=builder /app/server .

# Run
CMD ["./server"]