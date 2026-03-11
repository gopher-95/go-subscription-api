FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o subscription-api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscription-api .
COPY migrations ./migrations
COPY .env .env

EXPOSE 8080

CMD ["./subscription-api"]