FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./out/bin ./cmd/api/main.go

FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch

ENTRYPOINT []

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build ./app/out/bin ./out/bin

EXPOSE 80

ENTRYPOINT ["./out/bin"]

