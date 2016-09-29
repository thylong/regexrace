FROM golang:1.7.1

ADD . /go/src/github.com/thylong/regexrace
WORKDIR /go/src/github.com/thylong/regexrace

# Fetching Dependencies.
RUN go get github.com/tools/godep
RUN godep restore

RUN go install github.com/thylong/regexrace

EXPOSE 8080

CMD ["/go/bin/regexrace"]
