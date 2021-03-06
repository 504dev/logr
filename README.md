<div align="center">
  <a href="https://logr.info/demo">
    <img width="128" height="128" src="https://raw.githubusercontent.com/504dev/logr-front/master/static/logr.png">
  </a>
  <p>
    <b>logr</b> is an open source logger and metric service.
    <br>
    Helps you debug and analyze performance of your features.
    <br>
    Get to know your application better with <b>logr</b>.
  </p>
</div>

[Demo]

[![Logr](https://raw.githubusercontent.com/504dev/logr-front/master/static/preview.jpg)][Demo]

[Demo]: http://logr.info/demo

* logs looks like in your `Terminal`
* storing data in `ClickHouse`
* transport data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend
* Authorization by `GitHub`

## Usage
```javascript
const { Logr } = require('logr-node-client');

const conf = new Logr({
    udp: ':7776',
    publicKey: 'MCAwDQYJKoZIhvcNAQEBBQADDwAwDAIFAMg7IrMCAwEAAQ==',
    privateKey: 'MC0CAQACBQDIOyKzAgMBAAECBQCHaZwRAgMA0nkCAwDziwIDAL+xAgJMKwICGq0=',
});

const logr = conf.newLogger('hello.log');

logr.info('Hello, Logr!');
```

## Requirements
* Node.js `v12`
* Npm `v6`
* Golang `v1.13`
* ClickHouse `v20`
* Mysql `v5.7`


## Build & Run service
### Manual
1. Clone repository: \
    `git clone --recurse-submodules git@github.com:504dev/logr.git && cd logr`
2. Init config file:
    `make config`
3. Fill **config.yml**, see [Config](#config)
4. Build frontend:
    `make front`
5. Run:
    `make run`

### Docker-compose
1. Clone repository: \
    `git clone --recurse-submodules git@github.com:504dev/logr.git && cd logr`
2. Create **.env** file:
    `make env`
3. Set env variables in **.env** file
4. Run: ` docker-compose up -d`

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

* `client_id` and `client_secret` is Github App keys (optional. set empty, if not sure)
* `org` is organization restriction (if set, only org members can authorize)
* `jwt_secret` is random string (using to sign temporary authorization tokens)

## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)
* Python [github.com/504dev/logr-python-client](https://github.com/504dev/logr-python-client)
* PHP [github.com/504dev/logr-php-client](https://github.com/504dev/logr-php-client) (logger only, metrics not supported)

## Utils
* Watcher [github.com/504dev/logr-watch](https://github.com/504dev/logr-watch)

