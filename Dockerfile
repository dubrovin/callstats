FROM golang:alpine
RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . /go/src/app
COPY . /go/src/app/vendor/github.com/dubrovin/callstats

RUN go install -v ./...

CMD ["app"]

EXPOSE 8080