# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and other dependencies if needed
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy .env-template if needed or ensure your app handles missing .env via env vars
# COPY .env ./ # Usually good to NOT copy .env and use real environment variables in Render

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
