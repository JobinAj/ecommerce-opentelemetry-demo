FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Ensure dependencies are tidy
RUN go mod tidy

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
EXPOSE 8003

# Start the service
CMD ["/service"]


