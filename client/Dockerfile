FROM golang:1.18.3 AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY ./client ./client
COPY ./pow ./pow

RUN CGO_ENABLED=0 go build -o ./bin/client ./client/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin/client ./client

ENTRYPOINT ["./client"]