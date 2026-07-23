FROM golang:1.23-alpine AS builder

RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /migrate /usr/local/bin/
COPY --from=builder /app/api /usr/local/bin/
COPY --from=builder /app/migrations /migrations

WORKDIR /
EXPOSE 8080
CMD ["sh", "-c", "until migrate -path /migrations -database \"$DATABASE_URL\" up; do echo 'waiting for postgres...'; sleep 2; done && api"]
