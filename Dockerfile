# Stage 1: Build the Go binary
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

# Copy dependency files first and download dependencies for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Stage 2: Create a minimal production runner image
FROM alpine:3.19

WORKDIR /app

# Copy ca-certificates for secure external API requests (like Resend API)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary
COPY --from=builder /app/main .

# Expose the default container port
EXPOSE 8080

# Run the application
CMD ["./main"]
