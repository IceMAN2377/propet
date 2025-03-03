FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-service ./cmd/myapp

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-service .

# Runtime dependencies
RUN apk --no-cache add ca-certificates

EXPOSE 8080

ENV DB_CONN="postgres://user:password@db:5432/tasks?sslmode=disable"

CMD ["./task-service"]
