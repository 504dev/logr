FROM clickhouse/clickhouse-server:23.11.5-alpine as clickhouse
ENV CLICKHOUSE_DB=logrdb
ENV CLICKHOUSE_USER=logr
ENV CLICKHOUSE_PASSWORD=12345678

FROM clickhouse as mysql

RUN apk add --update --no-cache mysql mysql-client
COPY --chmod=644 ./Dockerhub/conf/mysql_startup.sh /mysql_startup.sh
COPY --chmod=644 ./Dockerhub/conf/my.cnf /etc/mysql/my.cnf
COPY --chmod=644 ./Dockerhub/conf/mariadb-server.cnf /etc/my.cnf.d/mariadb-server.cnf
EXPOSE 3306

FROM node:20.11-alpine AS frontend

COPY ./ /opt/logr
WORKDIR /opt/logr/frontend
RUN yarn install && yarn build && yarn cache clean

FROM golang:1.17-alpine AS gobuild

COPY --from=frontend /opt/logr /opt/logr
COPY ./Dockerhub/config.yml /opt/logr
WORKDIR /opt/logr
RUN go build -o logr-server ./cmd/server/main.go

FROM mysql as launcher

RUN apk add --no-cache git
COPY --from=gobuild /opt/logr /opt/logr
WORKDIR /opt/logr
EXPOSE 7776 7778
CMD sh /mysql_startup.sh && (/entrypoint.sh &) && ./logr-server --config="./config.yml" --retries=60

