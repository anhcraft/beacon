# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /beacon .

# --- Runtime stage ---
FROM alpine:3.22

# Install Docker CLI (needed for trigger commands like `docker stack deploy`)
RUN apk add --no-cache docker-cli

# Create a non-root user (socket permissions are handled via group mapping at runtime)
RUN addgroup -S beacon && adduser -S beacon -G beacon

WORKDIR /app

COPY --from=builder /beacon /app/beacon

# Default config path inside the container
# Mount your config at runtime: -v ./config.yaml:/app/config.yaml:ro
ENTRYPOINT ["/app/beacon"]
CMD ["-config", "/app/config.yaml"]
