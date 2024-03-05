## Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /order_system cmd/order_system/main.go cmd/order_system/wire_gen.go

## Run stage
FROM alpine

# Settings
RUN apk add -U tzdata ca-certificates
ENV TZ=America/Sao_Paulo
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime

COPY --from=builder /order_system /order_system
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/example.env /.env

ENTRYPOINT [ "/order_system" ]
