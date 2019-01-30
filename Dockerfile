FROM golang:1.11

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 50051
EXPOSE 5117

CMD ["server", "0.0.0.0:50051"]
