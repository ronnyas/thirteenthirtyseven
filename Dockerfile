# Builder stage
FROM golang:1.18-alpine AS builder
RUN apk add build-base
WORKDIR /app
COPY . .
RUN go build -o main

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# ENV TOKEN "token, if not using docker-compose"
# ENV MAINCHANNEL "..if not using docker-compose"
# ENV DATABASEPATH "...same"

CMD ["./main"]