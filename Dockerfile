# Builder
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o exprun ./cmd/main.go

# Runner
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata

COPY --from=builder /app/exprun .

CMD ["./exprun"]
