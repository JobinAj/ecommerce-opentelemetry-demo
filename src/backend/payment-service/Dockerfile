FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN go build -o service main.go 


# Use a minimal alpine image for the final image
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/service /service

# Expose the port the service runs on
EXPOSE 8003

# Start the service
CMD ["/service"]


