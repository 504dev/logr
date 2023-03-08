#!/usr/bin/env bash

shopt -s expand_aliases

DIR="$(dirname "${BASH_SOURCE[0]}")"

cd "${DIR}" && cd ..

if [ ! -f ".env" ]; then
    make env
fi

if [ "$(uname)" == "Darwin" ]; then
    alias sha256sum="shasum -a 256"
    alias sed="sed -i ''"
else
    alias sed="sed -i''"
fi

source .env

## Create clickhouse user from template
CLICKHOUSE_USER_FILE="./clickhouse/${CLICKHOUSE_USER}.xml"
CLICKHOUSE_PASSWORD_HASH=$(echo -n "${CLICKHOUSE_PASSWORD}" | sha256sum | awk '{ print $1 }')

cp templates/user.xml.template                                              "${CLICKHOUSE_USER_FILE}"

sed -e "s/\${CLICKHOUSE_USER\}/${CLICKHOUSE_USER}/g"                        "${CLICKHOUSE_USER_FILE}"
sed -e "s/\${CLICKHOUSE_PASSWORD_HASH\}/${CLICKHOUSE_PASSWORD_HASH}/g"      "${CLICKHOUSE_USER_FILE}"


CLICKHOUSE_INITDB_FILE="./clickhouse/init-db.sh"

cp templates/init-db.sh                                                     "${CLICKHOUSE_INITDB_FILE}"
sed -e "s/\${CLICKHOUSE_DATABASE\}/${CLICKHOUSE_DATABASE}/g"                "${CLICKHOUSE_INITDB_FILE}"


## Create Logr config from template
LOGR_CONFIG_FILE="config.yml"

cp templates/config.yml.template                                            "${LOGR_CONFIG_FILE}"

sed -e "s/\${LOGR_HTTP_HOST\}/${LOGR_HTTP_HOST}/g"                          "${LOGR_CONFIG_FILE}"
sed -e "s/\${LOGR_HTTP_PORT\}/${LOGR_HTTP_PORT}/g"                          "${LOGR_CONFIG_FILE}"
sed -e "s/\${LOGR_UDP_HOST\}/${LOGR_UDP_HOST}/g"                            "${LOGR_CONFIG_FILE}"
sed -e "s/\${LOGR_UDP_PORT\}/${LOGR_UDP_PORT}/g"                            "${LOGR_CONFIG_FILE}"

sed -e "s/\${MYSQL_HOST\}/${MYSQL_HOST}/g"                                  "${LOGR_CONFIG_FILE}"
sed -e "s/\${MYSQL_PORT\}/${MYSQL_PORT}/g"                                  "${LOGR_CONFIG_FILE}"
sed -e "s/\${MYSQL_ROOT_PASSWORD\}/${MYSQL_ROOT_PASSWORD}/g"                "${LOGR_CONFIG_FILE}"
sed -e "s/\${MYSQL_DATABASE\}/${MYSQL_DATABASE}/g"                          "${LOGR_CONFIG_FILE}"
sed -e "s/\${MYSQL_USER\}/${MYSQL_USER}/g"                                  "${LOGR_CONFIG_FILE}"
sed -e "s/\${MYSQL_PASSWORD\}/${MYSQL_PASSWORD}/g"                          "${LOGR_CONFIG_FILE}"

sed -e "s/\${CLICKHOUSE_HOST\}/${CLICKHOUSE_HOST}/g"                        "${LOGR_CONFIG_FILE}"
sed -e "s/\${CLICKHOUSE_PORT\}/${CLICKHOUSE_PORT}/g"                        "${LOGR_CONFIG_FILE}"
sed -e "s/\${CLICKHOUSE_DATABASE\}/${CLICKHOUSE_DATABASE}/g"                "${LOGR_CONFIG_FILE}"
sed -e "s/\${CLICKHOUSE_USER\}/${CLICKHOUSE_USER}/g"                        "${LOGR_CONFIG_FILE}"
sed -e "s/\${CLICKHOUSE_PASSWORD\}/${CLICKHOUSE_PASSWORD}/g"                "${LOGR_CONFIG_FILE}"

sed -e "s/\${OAUTH_JWT_SECRET\}/${OAUTH_JWT_SECRET}/g"                      "${LOGR_CONFIG_FILE}"
sed -e "s/\${OAUTH_GITHUB_ORG\}/${OAUTH_GITHUB_ORG}/g"                      "${LOGR_CONFIG_FILE}"
sed -e "s/\${OAUTH_GITHUB_CLIENT_ID\}/${OAUTH_GITHUB_CLIENT_ID}/g"          "${LOGR_CONFIG_FILE}"
sed -e "s/\${OAUTH_GITHUB_CLIENT_SECRET\}/${OAUTH_GITHUB_CLIENT_SECRET}/g"  "${LOGR_CONFIG_FILE}"