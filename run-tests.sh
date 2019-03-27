#!/bin/bash
set -ex
#
go get -u golang.org/x/lint/golint
go get -u github.com/jstemmer/go-junit-report
go get -u github.com/kardianos/govendor
govendor sync
golint -set_exit_status $(go list ./... | grep -v /vendor/) > golint-report.out
go vet $(go list ./... | grep -v /vendor/) 2> govet-report.out
go test -coverprofile=coverage.out -v 2>&1 | go-junit-report > report.xml

