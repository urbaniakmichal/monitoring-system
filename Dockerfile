# --- STAGE 1
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install swag tool
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

# Generate Swagger docs before building binary
RUN swag init -g cmd/monitor-agent/main.go --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o monitor-system ./cmd/monitor-agent

# --- STAGE 2
FROM debian:bookworm-slim

RUN groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/monitor-system .
COPY --from=builder --chown=appuser:appgroup /app/configs ./configs

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/monitor-system"]