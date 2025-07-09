FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
COPY --from=builder /app/cmd/migrations ./cmd/migrations
EXPOSE 8080
ENTRYPOINT ["./server"]
