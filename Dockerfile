# ── Stage 1: Build ──────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files first — cached layer, only re-downloads on go.mod changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 → statically linked binary (no glibc needed in runtime)
# -ldflags="-s -w" → strip debug symbols (~30% smaller binary)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ── Stage 2: Runtime ─────────────────────────────────────────────
FROM alpine:3.19

# ca-certificates needed for HTTPS calls (JWT library)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]