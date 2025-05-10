FROM golang:1.24.3-alpine3.21

WORKDIR /app

RUN apk add git

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o /coda-docker-linux main.go

EXPOSE 3000

CMD ["/coda-docker-linux", "server"]