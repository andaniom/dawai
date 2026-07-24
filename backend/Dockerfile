FROM golang:1.25-alpine AS builder

RUN go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY --from=builder /app/api /usr/local/bin/
COPY --from=builder /app/migrations /migrations

WORKDIR /
EXPOSE 8080
CMD ["sh", "-c", "until migrate -path /migrations -database \"$DATABASE_URL\" up; do echo 'waiting for postgres...'; sleep 2; done && api"]
