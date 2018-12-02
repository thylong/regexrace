FROM golang:1.11

ADD . /go/src/github.com/thylong/regexrace
WORKDIR /go/src/github.com/thylong/regexrace

# Fetching Dependencies.
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

RUN go install github.com/thylong/regexrace

EXPOSE 8080
CMD ["/go/bin/regexrace"]
