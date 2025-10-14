# Multi-stage Dockerfile for the Go backend (build -> small runtime)
# Build with static binary, run in minimal Alpine image with non-root user.

# ---------- builder ----------
FROM golang:1.25-alpine AS builder

# install tools needed for go modules and fetching
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# copy go mod files and download deps (cacheable)
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go build -ldflags="-s -w" -o /app/bin/web_chat ./main.go

# ---------- runtime ----------
FROM alpine:3.18

# create app user
RUN addgroup -S app && adduser -S -G app app

# required for TLS calls (e.g. Google, SMTP)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/bin/web_chat /app/web_chat
COPY .env.developer /app/.env
# optionally copy default env (do NOT store secrets in image for production)
# COPY .env /app/.env

RUN chown app:app /app/web_chat
USER app

EXPOSE 8080

# small healthcheck (optional)
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s CMD wget -qO- --timeout=2 http://127.0.0.1:8080/healthz >/dev/null || exit 1

ENTRYPOINT ["/app/web_chat"]