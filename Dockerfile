FROM golang:1.25.5-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o vtask ./cmd/vtask

FROM alpine:3.20

WORKDIR /app

RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/main" > /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.20/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache ca-certificates

COPY --from=builder /app/vtask .

EXPOSE 3003


CMD ["./vtask"]