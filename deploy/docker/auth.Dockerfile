# ---------- Build Stage ----------
FROM golang:1.25 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the auth binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o auth \
    ./cmd/auth


# ---------- Runtime Stage ----------
FROM debian:bookworm-slim

# Install CA certificates (needed for HTTPS requests, SMTP TLS, etc.)
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary
COPY --from=builder /app/auth .

# Copy web assets/templates (for SSR)
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./auth"]