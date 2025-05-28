# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.23-alpine AS builder
WORKDIR /app

RUN apk add --no-cache make

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/golang/mock/mockgen@latest && ln -s /go/bin/mockgen /usr/local/bin/mockgen

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make setup

RUN CGO_ENABLED=0 GOOS=linux go build -o gophernet ./cmd/main

# --- Run stage ---
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/gophernet ./gophernet
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/data ./data
COPY config.yaml /app/config.yaml

EXPOSE 8080

# Run the app
CMD ["./gophernet"] 