# ---------- Builder ----------
FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o cache-server ./cmd/cache-server

# ---------- Runtime ----------
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /src/cache-server .
COPY config/ ./config/
EXPOSE 8080
ENTRYPOINT ["./cache-server"]
