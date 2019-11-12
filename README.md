# CESSDA Café: Carsten's Coffee Machine

Golang implementation of CESSDA Café Coffee machine using [Gin Gonic](https://github.com/gin-gonic/gin/).

## Installation

Assuming you have go, and the code in you go path,
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

## Test suite

The test suite can be executed by calling

```bash
make test
```

## Documentation

To produce swagger.json, call
```bash
make swagger
```


## Execution

Simply run

```bash
make run
```

and access <http://localhost:1337.>

