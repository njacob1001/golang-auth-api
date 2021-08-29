FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./out/bin ./cmd/api/main.go

FROM scratch

COPY --from=build ./app/out/bin ./out/bin

EXPOSE 80

ENTRYPOINT ["./out/bin"]

