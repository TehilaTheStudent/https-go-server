# Start from the official Golang image for building
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go ./
COPY cert.pem ./
COPY key.pem ./
RUN go build -o server main.go

# Use a minimal base image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./
COPY --from=builder /app/cert.pem ./
COPY --from=builder /app/key.pem ./
EXPOSE 8443
CMD ["./server"]
