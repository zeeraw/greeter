FROM golang:1.11

WORKDIR /go/src/github.com/zeeraw/greeter
COPY . .
RUN go install -v ./client ./server

EXPOSE 50051
EXPOSE 5117

CMD ["server", "0.0.0.0:50051"]
