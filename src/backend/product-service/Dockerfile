# Use the official Golang alpine image as the base image
FROM golang:1.21-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY main.go .

# Build the Go application
RUN go build -o service main.go

# Use a minimal alpine image for the final image
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/service /service

# Expose the port the service runs on
EXPOSE 8001

# Start the service
CMD ["/service"]
