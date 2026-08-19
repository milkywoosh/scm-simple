# ---- Build stage ----
FROM golang:1.26.6-bookworm AS builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        git \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    .

# ---- Runtime stage ----
FROM alpine:3.19

# Install certificates and timezone data cleanly
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server .

# Ensure non-root user owns the binary (if needed) and switch user
USER app

EXPOSE 8000

ENTRYPOINT ["./server"]