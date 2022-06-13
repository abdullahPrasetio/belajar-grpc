FROM golang:1.17-alpine

WORKDIR /app


RUN apk --no-cache add curl


COPY . .

RUN go get

RUN go build -o belajar-grpc

EXPOSE 8080


CMD ./belajar-grpc