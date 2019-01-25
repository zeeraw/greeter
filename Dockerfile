FROM golang:1.11

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server", "0.0.0.0:50051"]
