# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o bot ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache ca-certificates sqlite

# Copy binary from builder
COPY --from=builder /app/bot .

# Create directories for persistent data
RUN mkdir -p /app/session /app/credentials

# Set environment variables
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/credentials.json
ENV WHATSAPP_SESSION_DB_PATH=/app/session/session.db

# Run the bot
CMD ["./bot"]
