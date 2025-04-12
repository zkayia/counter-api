# Deps stage
FROM golang:1.24-alpine AS deps

WORKDIR /app

COPY counter-api/go.mod counter-api/go.sum ./

RUN go mod download

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY counter-api/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN addgroup -g 1000 counterapi
RUN adduser -D -s /bin/sh -u 1000 -G counterapi counterapi

COPY --from=builder /app/main .

RUN chown counterapi:counterapi /app/main

USER counterapi
EXPOSE 8000
ENTRYPOINT ["/app/main"]
