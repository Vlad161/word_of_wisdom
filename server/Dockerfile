FROM golang:1.18.3 AS build

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY ./server /src/server
COPY ./pow /src/pow

RUN CGO_ENABLED=0 go build -o ./bin/server ./server/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=build /src/bin/server /app/server

ENTRYPOINT ["/app/server"]