FROM golang:latest

COPY . /go/src/github.com/vektorlab/otter
RUN cd /go/src/github.com/vektorlab/otter && \
    go get -v -d

RUN go install github.com/vektorlab/otter
