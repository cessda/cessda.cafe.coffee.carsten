# CESSDA Café: Carsten's Coffee Machine

Golang implementation of CESSDA Café Coffee machine using [Gin Gonic](https://github.com/gin-gonic/gin/).

## Installation

Assuming you have go, and the code in your go path,
install the dependencies with [Govendor](https://github.com/kardianos/govendor).

```bash
govendor sync
```

## Options

To specify a port other then the default, set the environment variable
```bash
export COFFEE_PORT="1337"
```

To log log in Gelf format, use the follwing environment setting
```bash
export GELF_LOGGING="true"
```

In order to turn off debug mode, use
```bash
export GIN_MODE="release"
```

## Test suite

The test suite can be executed by calling

```bash
make test
```

## API Spec

To produce swagger.json, call
```bash
make swagger
```


## Execution

Simply run

```bash
make run
```

and access <http://localhost:1337>.

## Use

Consider the following call:
```bash
curl -s -XPOST -d '{ "product": "COFFEE", "jobId": "00000000-BBBB-0000-BBBB-000000000000" }' \
   -H "Content-Type: application/json" http://localhost:1337/start-job
```
with response
```json
{
  "jobId": "00000000-BBBB-0000-BBBB-000000000000",
  "product": "COFFEE",
  "jobStarted": "2019-08-01T00:00:01.000Z",
  "jobReady": "2019-08-01T00:00:31.000Z"
}
```

Follow it up with
```bash
curl -s http://localhost:1337/retrieve-job/00000000-BBBB-0000-BBBB-000000000000
```
to get
```json
{
  "message": "Job not ready"
}
```
and later
```json
{
  "jobId": "00000000-BBBB-0000-BBBB-000000000000",
  "product": "COFFEE",
  "jobStarted": "2019-08-01T00:00:01.000Z",
  "jobReady": "2019-08-01T00:00:31.000Z",
  "jobRetrieved": "2019-08-01T00:01:31.000Z"
}
```

