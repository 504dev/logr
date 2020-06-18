# Logr

_Logr_ is an open source logger and counter. \
Get to know your application better.

* logs looks like in your `Terminal`
* counters are drawn with `Highcharts`
* storing data in `ClickHouse`
* updating data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend

## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)

## Build service
1. Create directories:
    `cd $GOPATH/src/github.com && mkdir 504dev && cd $_`
2. Clone repository:
    `git clone git@github.com:504dev/logr.git && cd logr`
3. Make helper:
    `make`
4. Init config file:
    `make config`
5. Build & run:
    `make run`