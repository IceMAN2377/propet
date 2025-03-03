FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-service .

EXPOSE 8000

CMD ["./api-service"]
