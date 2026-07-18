# ============================================================
# Stage 1: Build
# ============================================================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/console ./cmd/console

# ============================================================
# Stage 2: Run
# ============================================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/console .
COPY --from=builder /app/.env .env 2>/dev/null || true
COPY --from=builder /app/templates ./templates 2>/dev/null || true
COPY --from=builder /app/plugins ./plugins 2>/dev/null || true

EXPOSE ${APP_PORT:-8200}

CMD ["./server"]
