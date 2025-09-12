# Build stage
FROM golang:1.21-alpine AS builder

# Install ca-certificates for SSL/TLS
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Final stage
FROM alpine:latest

# Install ca-certificates for SSL/TLS
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8080

CMD ["./main"]
