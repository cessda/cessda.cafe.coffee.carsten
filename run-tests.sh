#!/bin/bash

# Carsten's CESSDA CAFE Coffee Machine
# Copyright CESSDA-ERIC 2019
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not
# use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex
#
go get -u golang.org/x/lint/golint
go get -u github.com/jstemmer/go-junit-report
go get -u github.com/kardianos/govendor
govendor sync
# linting and standard code vetting
golint -set_exit_status $(go list ./... | grep -v /vendor/) > golint-report.out
go vet $(go list ./... | grep -v /vendor/) 2> govet-report.out
# test report in junit-format
go test -v 2>&1 | go-junit-report > junit.xml
# coverage report for use by sonar
go test -coverprofile=coverage.out -v -json 2>&1 > report.json

