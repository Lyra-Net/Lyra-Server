# ---- Build stage ----
FROM golang:1.25.1-alpine AS builder

WORKDIR /app


# Copy go mod và download dependency trước để cache
COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

# Copy toàn bộ code
COPY . .

# Build binary
RUN go build -o search-service main.go

# ---- Run stage ----
FROM alpine:3.20

WORKDIR /app


COPY --from=builder /app/search-service .


COPY .env .

EXPOSE 30005

CMD ["./search-service"]
