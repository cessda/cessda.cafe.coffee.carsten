#vetting and running the test suite:
go get -u github.com/jstemmer/go-junit-report
go get -u github.com/kardianos/govendor
govendor sync
go vet $(go list ./... | grep -v /vendor/)
go test -v 2>&1 | go-junit-report > report.xml