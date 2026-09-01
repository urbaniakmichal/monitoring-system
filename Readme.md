# 🚀 Monitor System - Backend Microservice

The core backend service for the **Monitor System** infrastructure, built with **Go**. It provides REST API endpoints, business logic, structured logging, Elasticsearch integration, and native Prometheus metrics exposition.

---

## 📋 Table of Contents
- [About the Project](#-about-the-project)
- [Tech Stack](#-tech-stack)
- [Prerequisites](#-prerequisites)
- [Project Structure](#-project-structure)
- [Environment Variables](#-environment-variables)
- [Getting Started](#-getting-started)
- [API Endpoints & Metrics](#-api-endpoints--metrics)

---

## 🎯 About the Project

`monitor-system` is designed for high performance and reliability. It exposes system health, logs analytics via Elasticsearch, and serves internal telemetry metrics (`/metrics`) to be scraped by Prometheus and visualized in Grafana.

---

## 🛠️ Tech Stack

* **Language:** Go
* **Metrics:** Prometheus (`promhttp`)
* **Logging & Storage:** Elasticsearch Integration
* **Containerization:** Docker & Docker Compose

---

## 📦 Prerequisites

Ensure you have the following installed on your machine:
* [Go](https://golang.org/) (version 1.21+ recommended)
* [Docker & Docker Compose](https://www.docker.com/) (for running the full stack)

---

## 📂 Project Structure

```text
monitor-system/
├── cmd/                  # Application entrypoints
├── internal/             # Private application code (handlers, services)
├── configs/              # Configuration files (Prometheus, Grafana)
├── Dockerfile            # Multi-stage Docker build file
└── README.md             # Project documentation
```

---

## ⚙️ Environment Variables

The application can be configured using the following environment variables:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Port on which the HTTP server listens | `8080` |
| `ELASTICSEARCH_URL` | URL connection string for Elasticsearch | `http://elasticsearch:9200` |

---

## 🚀 Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/urbaniakmichal/monitor-system.git
   cd monitor-system
   ```

2. **Run the application locally:**
   ```bash
   go run cmd/main.go
   ```

3. **Or run via Docker Compose (with the full monitoring stack):**
   ```bash
   docker compose up --build
   ```

---

## 🔌 API Endpoints & Metrics

* **Health Check:** `GET /api/v1/health` – Returns service status.
* **Prometheus Metrics:** `GET /metrics` – Exposes application telemetry (goroutines, CPU, memory, HTTP latency counters).

---

## 📄 License

This project is licensed under the MIT License.