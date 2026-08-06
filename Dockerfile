# Stage 1: Build the binary 
FROM golang:1.26.5 AS builder
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod ./
RUN go mod download

# Copy SRC files and build the Go binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Final runtime Alpine image
FROM scratch

# Copy binary from the builder stage
COPY --from=builder /app/server .

# Run
CMD ["./server"]