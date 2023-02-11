FROM golang:1.19-alpine3.16 AS builder

# Create and change to the /build directory
WORKDIR /build

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy everything to the /build directory
COPY . .

# Build book-manager
RUN go build -o ptask cmd/main.go

# Deploy
FROM alpine:latest

# Update and add timezone data
RUN apk update && apk add bash && apk --no-cache add tzdata

# Change to the /build directory
WORKDIR /build

# Copy everything from builder to the /build directory
COPY --from=builder /build/ptask .

ENTRYPOINT ["./ptask"]