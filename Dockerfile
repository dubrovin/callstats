FROM golang:alpine
RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . /go/src/app
COPY ./controller /go/src/app/vendor/github.com/dubrovin/callstats/controller
COPY ./server /go/src/app/vendor/github.com/dubrovin/callstats/server

RUN go install -v ./...

CMD ["app", "--delay=1s"]

EXPOSE 8080