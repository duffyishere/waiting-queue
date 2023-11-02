FROM golang:1.21 AS build

WORKDIR /build

COPY core/go.mod ./
COPY core/go.sum ./
COPY core/*.go ./

RUN go mod download && go mod verify

RUN go build -o app main.go redisUtils.go

EXPOSE 3000

ENTRYPOINT ["./app"]