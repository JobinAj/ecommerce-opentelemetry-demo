# Use the official Golang alpine image as the base image
FROM golang:1.24-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Ensure dependencies are tidy
RUN go mod tidy

# Build the Go application
RUN go build -o service .

# Use a minimal alpine image for the final image
FROM alpine:latest

# Install ca-certificates and created a non-root user
RUN apk --no-cache add ca-certificates && \
    addgroup -S ecom && \
    adduser -S ecom -G ecom -u 1000

# Copy the binary from the builder stage
COPY --from=builder /app/service /service
RUN chown ecom:ecom /service

# Use the non-root user
USER ecom

# Expose the port the service runs on
EXPOSE 8002

# Start the service

CMD ["/service"]
