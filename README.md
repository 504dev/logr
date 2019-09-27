# kidlog
logger

0) install modules gin, sqlx ... (using go module system)
1) clickhouse.Init() init conn, create schemas
2) read log from clickhouse, use sqlx
3) http and udp handlers
4) write log to clickhouse

# Build service

1. Create directories:
    `cd $GOPATH/src/github.com && mkdir 504dev && cd $_`
2. Clone repository:
    `git clone git@github.com:504dev/kidlog.git && cd kidlog`
3. Make helper:
    `make`
4. Init config file:
    `make config`
5. Build & run:
    `make run`