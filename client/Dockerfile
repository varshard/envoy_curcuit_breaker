FROM golang:latest

WORKDIR /go/src/github.com/varshard/client/

EXPOSE 6061

COPY main.go .

RUN go build -o main .

CMD ["./main"]
