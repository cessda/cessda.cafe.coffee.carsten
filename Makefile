# COPYRIGHT CESSDA ERIC 2019
# SPDX-Identifier: Apache-2.0
#


prep:
	go get -u golang.org/x/lint/golint
	go get -u github.com/jstemmer/go-junit-report
	go get -u github.com/kardianos/govendor
	govendor sync


lint:
	# linting and standard code vetting
	golint -set_exit_status $(go list ./... | grep -v /vendor/) | tee golint-report.out
	go vet $(go list ./... | grep -v /vendor/) 2> govet-report.out


test: prep lint
	# presentable output
	go test -v

test-ci: prep lint
	# test report in junit-format
	go test -v 2>&1 | go-junit-report > junit.xml
	# coverage report for use by sonar
	go test -coverprofile=coverage.out -v -json 2>&1 > report.json

build: prep
	go build -v .

run:
	go run .

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.json --scan-models

