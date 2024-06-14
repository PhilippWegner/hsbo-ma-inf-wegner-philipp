FROM golang:11.22.4-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o statemachine-service-app ./cmd/api

RUN chmod +x /app/statemachine-service-app

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/statemachine-service-app /app

CMD [ "/app/statemachine-service-app" ]
