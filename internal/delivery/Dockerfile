FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /app

COPY internal/delivery/. delivery/
COPY . .

RUN go mod tidy
RUN go build -o main delivery/main.go

FROM alpine:edge AS prod

RUN apk add bash

COPY --from=build /app/main /app/main
COPY --from=build /app/.env /app/.env

RUN mkdir -p /app/configs
COPY configs app/configs

WORKDIR /app
EXPOSE 8000

ENTRYPOINT ["./main"]