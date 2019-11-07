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

To send logs in Gelf format to a Graylog server via UDP, provide the address and port
```bash
export GRAYLOG_SERVER="<server-ip>:<server-port>"
```


## Execution

Simply run

```bash
go run .
```

and access <http://localhost:1337.>

