FROM golang:1.22.4-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o restful-broker-service-app ./cmd/api

RUN chmod +x /app/restful-broker-service-app

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/restful-broker-service-app /app

CMD [ "/app/restful-broker-service-app" ]
