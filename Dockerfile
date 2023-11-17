FROM golang:alpine
WORKDIR $GOPATH/data
COPY . $GOPATH/data
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux GOARCH=amd64 go build main.go
VOLUME ["/data/config","/data/log"]
ENTRYPOINT ["./main"]

