# Build stage
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o taskkr ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/taskkr .

COPY .env .env

EXPOSE 8080

CMD ["./taskkr"]
