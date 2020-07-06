# Logr

![Hello World](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEYAAAAUCAAAAAAVAxSkAAABrUlEQVQ4y+3TPUvDQBgH8OdDOGa+oUMgk2MpdHIIgpSUiqC0OKirgxYX8QVFRQRpBRF8KShqLbgIYkUEteCgFVuqUEVxEIkvJFhae3m8S2KbSkcFBw9yHP88+eXucgH8kQZ/jSm4VDaIy9RKCpKac9NKgU4uEJNwhHhK3qvPBVO8rxRWmFXPF+NSM1KVMbwriAMwhDgVcrxeMZm85GR0PhvGJAAmyozJsbsxgNEir4iEjIK0SYqGd8sOR3rJAGN2BCEkOxhxMhpd8Mk0CXtZacxi1hr20mI/rzgnxayoidevcGuHXTC/q6QuYSMt1jC+gBIiMg12v2vb5NlklChiWnhmFZpwvxDGzuUzV8kOg+N8UUvNBp64vy9q3UN7gDXhwWLY2nMC3zRDibfsY7wjEkY79CdMZhrxSqqzxf4ZRPXwzWJirMicDa5KwiPeARygHXKNMQHEy3rMopDR20XNZGbJzUtrwDC/KshlLDWyqdmhxZzCsdYmf2fWZPoxCEDyfIvdtNQH0PRkH6Q51g8rFO3Qzxh2LbItcDCOpmuOsV7ntNaERe3v/lP/zO8yn4N+yNPrekmPAAAAAElFTkSuQmCC)


_Logr_ is an open source logger and counter. \
Get to know your application better.

* logs looks like in your `Terminal`
* storing data in `ClickHouse`
* transport data by `WebSocket`
* `Golang` backend
* `Vue.js` frontend

## Requirements
* Node.js `v12`
* Npm `v6`
* Golang `v1.13`
* ClickHouse `v20`
* Mysql `v5.7`

## Build service
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
    client_id: "client_id"
    client_secret: "client_secret"
  jwt_secret: "jwt-secret"
clickhouse: "tcp://localhost:9000?database=logr&username=logr&password=logr"
mysql: "logr:logr@/logr"
```

* Create Github OAuth App, set `client_id` and `client_secret` \
https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/
* `jwt_secret` is random string (using to sign temporary authorization token)


## Client libraries

* Golang [github.com/504dev/logr-go-client](https://github.com/504dev/logr-go-client)
* Node.js [github.com/504dev/logr-node-client](https://github.com/504dev/logr-node-client)

## Utils
* Watcher [github.com/504dev/logr-watch](https://github.com/504dev/logr-watch)

