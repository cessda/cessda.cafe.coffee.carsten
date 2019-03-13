#before_script:
mkdir -p /go/src/nsd-utvikling /go/src/_/builds
cp -r . /go/src/nsd-utvikling/carsten-coffee-api
ln -s /go/src/nsd-utvikling /go/src/_/builds/nsd-utvikling

#golang linting:
go get -u golang.org/x/lint/golint
golint -set_exit_status $(go list ./... | grep -v /vendor/)

#vetting and running the test suite:
cd /go/src/nsd-utvikling/carsten-coffee-api
go get -u github.com/jstemmer/go-junit-report
go get -u github.com/kardianos/govendor
govendor sync
go vet $(go list ./... | grep -v /vendor/)
go test -v 2>&1 | go-junit-report > report.xml