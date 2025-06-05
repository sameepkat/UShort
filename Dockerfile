# Build stage
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ushort ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ushort .

EXPOSE 8080

CMD ["./ushort"]