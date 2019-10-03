FROM golang:latest

WORKDIR /go/src/github.com/varshard/pingpong/

EXPOSE 6060

COPY main.go .

RUN go build -o main .

CMD ["./main"]
