FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/app

FROM alpine:latest
COPY --from=builder /app/avito-backend-intership-2025 /usr/local/bin/avito-backend-intership-2025
ENTRYPOINT ["server"]
