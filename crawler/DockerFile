# Use the official Golang image to build the binary
FROM golang:1.23 AS builder

# Set working directory in the container
WORKDIR /app

# Copy Go source files
COPY .. .

# Build the Go binary
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s' -o main .


FROM debian:stable

COPY --from=builder /app/main ./main

CMD [ "./main" ]