#!/bin/bash
set -ex
#
go get -u golang.org/x/lint/golint
go get -u github.com/jstemmer/go-junit-report
go get -u github.com/kardianos/govendor
govendor sync
golint -set_exit_status $(go list ./... | grep -v /vendor/)
go vet $(go list ./... | grep -v /vendor/)
go test -coverprofile=coverage.out -v 2>&1 | go-junit-report > report.xml
