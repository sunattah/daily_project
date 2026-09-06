# --- Build stage ---
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o daybook .

# --- Run stage ---
FROM debian:bookworm-slim

WORKDIR /app

# ca-certificates needed if the app ever makes outbound HTTPS calls
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/daybook .
COPY --from=builder /app/templates ./templates

ENV DB_PATH=/data/routine.db

EXPOSE 8080

CMD ["./daybook"]
