# Logr

_Logr_ is an open source logger and counter. \
Get to know your application better.

* logs looks like in your `Terminal`
* storing data in `ClickHouse`
* transport data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend

## Build service
1. Clone repository:
    `git clone git@github.com:504dev/logr.git && cd logr`
2. Init config file:
    `make config`
3. Fill config, see Config part
4. Build frontend:
    `make front`
5. Run:
    `make run`

## Config
```yaml
bind:
  http: ":7778"
  udp: ":7776"
oauth:
  github:
    client_id: "client_id"
    client_secret: "client_secret"
  jwt_secret: "jwt-secret"
clickhouse: "tcp://localhost:9000?database=logr&username=logr&password=logr"
mysql: "logr:logr@/logr"
```

## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)

## Utils
* Watcher [github.com/504dev/logr-watch](https://github.com/504dev/logr-watch)

