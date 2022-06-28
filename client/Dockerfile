FROM golang:1.18.3 AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY ./client /src/client
COPY ./pow /src/pow

RUN CGO_ENABLED=0 go build -o ./bin/client ./client/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin/client /app/client

ENTRYPOINT ["/app/client"]