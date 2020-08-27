#!/usr/bin/env bash

DIR="$(dirname "${BASH_SOURCE[0]}")" 

cd ${DIR} && cd ..

if [ ! -f ".env" ]; then
    echo "Error! You need to create .env file from .env.example first!"
    exit 1
fi

if ! sha256sum; then
    function sha256sum() { shasum -a 256 "$@" ; } && export -f sha256sum
fi

source .env

## Create clickhouse user from template
CLICKHOUSE_USER_FILE="${CLICKHOUSE_USER}.xml"
CLICKHOUSE_PASSWORD_HASH=$(echo -n "${CLICKHOUSE_PASSWORD}" | sha256sum | awk '{ print $1 }')

cp templates/user.xml.template "${CLICKHOUSE_USER_FILE}"

sed -i'' -e "s/\${CLICKHOUSE_USER\}/${CLICKHOUSE_USER}/g" "${CLICKHOUSE_USER_FILE}"
sed -i'' -e "s/\${CLICKHOUSE_PASSWORD_HASH\}/${CLICKHOUSE_PASSWORD_HASH}/g" "${CLICKHOUSE_USER_FILE}"

## Create Logr config from template
LOGR_CONFIG_FILE="config.yml"

cp templates/config.yml.template "${LOGR_CONFIG_FILE}"

sed -i'' -e "s/\${LOGR_HTTP_BIND\}/${LOGR_HTTP_BIND}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${LOGR_UDP_BIND\}/${LOGR_UDP_BIND}/g" "${LOGR_CONFIG_FILE}"

sed -i'' -e "s/\${MYSQL_HOST\}/${MYSQL_HOST}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${MYSQL_PORT\}/${MYSQL_PORT}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${MYSQL_ROOT_PASSWORD\}/${MYSQL_ROOT_PASSWORD}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${MYSQL_DATABASE\}/${MYSQL_DATABASE}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${MYSQL_USER\}/${MYSQL_USER}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${MYSQL_PASSWORD\}/${MYSQL_PASSWORD}/g" "${LOGR_CONFIG_FILE}"

sed -i'' -e "s/\${CLICKHOUSE_HOST\}/${CLICKHOUSE_HOST}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${CLICKHOUSE_PORT\}/${CLICKHOUSE_PORT}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${CLICKHOUSE_DATABASE\}/${CLICKHOUSE_DATABASE}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${CLICKHOUSE_USER\}/${CLICKHOUSE_USER}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${CLICKHOUSE_PASSWORD\}/${CLICKHOUSE_PASSWORD}/g" "${LOGR_CONFIG_FILE}"

sed -i'' -e "s/\${OAUTH_JWT_SECRET\}/${OAUTH_JWT_SECRET}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${OAUTH_GITHUB_ORG\}/${OAUTH_GITHUB_ORG}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${OAUTH_GITHUB_CLIENT_ID\}/${OAUTH_GITHUB_CLIENT_ID}/g" "${LOGR_CONFIG_FILE}"
sed -i'' -e "s/\${OAUTH_GITHUB_CLIENT_SECRET\}/${OAUTH_GITHUB_CLIENT_SECRET}/g" "${LOGR_CONFIG_FILE}"

