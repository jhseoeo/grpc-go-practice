FROM golang:1.18

RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ADD . /app
WORKDIR /app

RUN go mod tidy && \
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pb/server_streaming.proto && \
    go build -o /app/server github.com/junhyuk0801/golang-grpc-practice/server_streaming/server

EXPOSE 8484

CMD /app/server