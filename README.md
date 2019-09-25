# kidlog
logger

0) install modules gin, sqlx ... (using go module system)
1) clickhouse.Init() init conn, create schemas
2) read log from clickhouse, use sqlx
3) http and udp handlers
4) write log to clickhouse
