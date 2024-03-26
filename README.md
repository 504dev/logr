<div align="center">
  <a href="https://logr.info/demo">
    <img width="128" height="128" src="https://raw.githubusercontent.com/504dev/logr-front/master/static/logr.png">
  </a>
  <p>
    <b>logr</b> is an open-source logger and metric service.
    <br>
    Helps you debug and analyze performance of your features.
    <br>
    Get to know your application better with <b>logr</b>.
  </p>
</div>

[![Logr](https://raw.githubusercontent.com/504dev/logr-front/master/static/preview.png)](http://logr.info/demo)

[Look at the demo page](http://logr.info/demo)


* logs look like in your `Terminal`
* storing data in `ClickHouse`
* transport data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend
* Authorization by `GitHub`

## Usage

```bash
docker run -d -p 7776:7776/udp -p 7778:7778 --name logr kozhurkin/logr
```

```javascript
const { Logr } = require('logr-node-client');

const conf = new Logr({
  udp: ':7776',
  publicKey: 'MCAwDQYJKoZIhvcNAQEBBQADDwAwDAIFAMg7IrMCAwEAAQ==',
  privateKey: 'MC0CAQACBQDIOyKzAgMBAAECBQCHaZwRAgMA0nkCAwDziwIDAL+xAgJMKwICGq0=',
});

const logr = conf.newLogger('hello.log');

logr.info('Hello, Logr!');

// 2024-01-30T22:50:04+03:00 info [v1.0.41, pid=60512, cmd/hello.go:41] Hello, Logr!
```

## Build & Run service

### Docker image

```
docker run -d -p 7776:7776/udp -p 7778:7778 --name logr kozhurkin/logr
```

Enjoy: \
http://localhost:7778/

### Docker-compose

1. Clone repository: \
   `git clone --recurse-submodules https://github.com/504dev/logr.git && cd logr`
2. Generate **.env** file: \
   `make env`
3. Edit **.env** file with your favourite editor
4. Generate **config.yml**: \
   `make config`
5. Run: \
   `docker-compose up -d`
6. Enjoy: \
   http://localhost:7778/

### Manual


Requirements: `Node.js v20` `Npm v10` `Golang v1.19` `ClickHouse v23` `Mysql v5.7`

1. Clone repository: \
   `git clone --recurse-submodules https://github.com/504dev/logr.git && cd logr`
2. Init **config.yml** file: \
   `make config`
3. Fill **config.yml**, see [Config](#config) section
4. Creating databases in Clickhouse and Mysql:
    ```
    clickhouse-client --query "CREATE DATABASE IF NOT EXISTS logrdb"
    mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS logrdb;"
    ```
5. Build frontend: \
   `make front`
6. Build backend: \
   `make build`
7. Run: \
   `make run`
8. Enjoy: \
   http://localhost:7778/

## Config

```yaml
bind:
  http: ":7778"
  udp: ":7776"
oauth:
  jwt_secret: "santaclausdoesntexist"
  github:
    client_id: "9bd30997b0ee30997b0ee3"
    client_secret: "1f241d37d910b11f241d37d910b11f241d37d910b1"
    org: "504dev"
clickhouse: "tcp://localhost:9000?database=logr&username=logr&password=logr"
mysql: "logr:logr@tcp(localhost:3306)/logr"
```

* `jwt_secret` is random string (using to sign temporary authorization tokens)
* `client_id` and `client_secret` is GitHub App keys (optional. set empty, if not sure)
* `org` is organization restriction (if set, only org members can authorize)

## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)
* Python [github.com/504dev/logr-python-client](https://github.com/504dev/logr-python-client)
* PHP [github.com/504dev/logr-php-client](https://github.com/504dev/logr-php-client) (logger only, metrics not
  supported)

## Utils

* Watcher [github.com/504dev/logr-watch](https://github.com/504dev/logr-watch)

