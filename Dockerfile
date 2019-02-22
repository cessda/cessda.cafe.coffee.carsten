# build the binary
FROM golang:latest as builder
WORKDIR /go/src/carsten-coffee-api
COPY *.go /go/src/carsten-coffee-api/
COPY vendor/vendor.json /go/src/carsten-coffee-api/vendor/
RUN go get -u github.com/kardianos/govendor
RUN govendor sync
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v .

# package the binary into a container
FROM scratch
COPY --from=builder /go/src/carsten-coffee-api/carsten-coffee-api /carsten-coffee-api
CMD ["/carsten-coffee-api"]


