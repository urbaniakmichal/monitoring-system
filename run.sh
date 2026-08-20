#!/usr/bin/env bash
set -e

swag init -g cmd/monitor-agent/main.go --parseDependency --parseInternal

go run ./cmd/monitor-agent/main.go