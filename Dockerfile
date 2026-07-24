# --- STAGE 1: Budowanie binarki ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o monitor-system ./cmd/monitor-agent

# --- STAGE 2: Obraz produkcyjny ---
FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# 1. Kopiujemy binarkę z builder
COPY --from=builder /app/monitor-system .

# 2. Kopiujemy folder z konfiguracją
COPY configs/ ./configs/

# Zmieniamy właściciela plików na appuser
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/monitor-system"]