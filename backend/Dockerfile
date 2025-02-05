# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

# Use a smaller runtime image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Set the working directory
WORKDIR /app

# Expose the application port
EXPOSE 8080

# Define the command to run the application
CMD ["main"]