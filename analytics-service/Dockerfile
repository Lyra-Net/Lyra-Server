# ---- Build stage ----
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy go mod & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN go install -tags='no_postgres no_mysql no_sqlite3 no_ydb no_mssql no_vertica no_libsql' \
    github.com/pressly/goose/v3/cmd/goose@latest

# Build binary service
RUN go build -o analytics-service cmd/main.go


# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates bash

# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration scripts
COPY --from=builder /app/db/migrations ./migrations

# Copy entrypoint script
COPY entrypoint.sh .

# Copy compiled service binary
COPY --from=builder /app/analytics-service .

COPY .env .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
