# ---- Build stage ----
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy go mod & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build binary service
RUN go build -o auth-service main.go

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app

# Install postgres client
RUN apk add --no-cache postgresql-client ca-certificates

# Copy migrate binary
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy migration scripts
COPY --from=builder /app/db/migration ./migration

# Copy templates
COPY --from=builder /app/templates ./templates

# Copy entrypoint script
COPY entrypoint.sh .

# Copy compiled service binary
COPY --from=builder /app/auth-service .

COPY .env .

EXPOSE 30003

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
