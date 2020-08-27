# Logr

_Logr_ is an open source logger and metric service. \
Get to know your application better.

[Demo](http://212.224.113.196:7778/demo)

![Logr](https://i.ibb.co/4dsbDdk/image.png)

* logs looks like in your `Terminal`
* storing data in `ClickHouse`
* transport data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend
* Authorization by `GitHub`

## Requirements
* Node.js `v12`
* Npm `v6`
* Golang `v1.13`
* ClickHouse `v20`
* Mysql `v5.7`

## Build & Run service
1. Clone repository: \
    `git clone --recurse-submodules git@github.com:504dev/logr.git && cd logr`
2. Init config file:
    `make config`
3. Fill config, see **Config**
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
    client_id: "9bd30997b0ee30997b0ee3"
    client_secret: "1f241d37d910b11f241d37d910b11f241d37d910b1"
    org: "504dev"
  jwt_secret: "jwt-secret"
clickhouse: "tcp://localhost:9000?database=logr&username=logr&password=logr"
mysql: "logr:logr@/logr"
```

* `client_id` and `client_secret`, need to create Github OAuth App \
https://docs.github.com/en/developers/apps/creating-an-oauth-app/
* `org` is organization restriction (if set, only members can authorize)
* `jwt_secret` is random string (using to sign temporary authorization tokens)

## Run in Docker
1. Clone repository: \
    `git clone --recurse-submodules git@github.com:504dev/logr.git && cd logr`
2. Create **.env** file:
    `make env`
3. Fill **.env** file
4. Run: ` docker-compose up -d`


## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)
* Python [github.com/504dev/logr-python-client](https://github.com/504dev/logr-python-client)

## Utils
* Watcher [github.com/504dev/logr-watch](https://github.com/504dev/logr-watch)

