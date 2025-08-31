# ---- Build stage ----
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download

COPY . .

RUN go build -o song-service main.go

# ---- Run stage ----
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/song-service .

COPY .env .

EXPOSE 30004

CMD ["./song-service"]
