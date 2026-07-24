# --- STAGE 1
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o monitor-system ./cmd/monitor-agent

# --- STAGE 2
FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/monitor-system .
COPY --from=builder --chown=appuser:appgroup /app/configs ./configs

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/monitor-system"]