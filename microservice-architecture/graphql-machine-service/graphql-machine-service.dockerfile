FROM golang:11.22.4-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o graphql-machine-service-app ./cmd/api

RUN chmod +x /app/graphql-machine-service-app

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/graphql-machine-service-app /app

CMD [ "/app/graphql-machine-service-app" ]
