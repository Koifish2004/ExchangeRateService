FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exchange-rate-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates


WORKDIR /root/

COPY --from=builder /app/exchange-rate-service .

EXPOSE 8080

CMD ["./exchange-rate-service"]
