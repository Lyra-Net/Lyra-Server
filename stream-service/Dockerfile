# ---- Build stage ----
FROM golang:1.23.5-alpine AS builder
WORKDIR /app

# Copy go mod & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .


# Build binary service
RUN go build -o stream-service main.go

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app

# Copy compiled service binary
COPY --from=builder /app/stream-service .
COPY .env .

EXPOSE 3002

CMD ["./stream-service"]
