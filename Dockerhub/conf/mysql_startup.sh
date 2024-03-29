#!/bin/sh

if [ ! -d "/run/mysqld" ]; then
  mkdir -p /run/mysqld
fi

if [ ! -z "$(ls -A /var/lib/mysql)" ]; then
  echo "[i] MySQL directory already present, skipping creation"
  /usr/bin/mysqld &
else
  echo "[i] MySQL data directory not found, creating initial DBs"

  mysql_install_db --user=root > /dev/null

  MYSQL_DATABASE=${MYSQL_DATABASE:-"logrdb"}
  MYSQL_USER=${MYSQL_USER:-"logr"}
  MYSQL_PASSWORD=${MYSQL_PASSWORD:-"12345678"}
  MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-"12345678"}

  tfile=`mktemp`
  if [ ! -f "$tfile" ]; then
      return 1
  fi

  cat << EOF > $tfile
USE mysql;
FLUSH PRIVILEGES;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY "$MYSQL_ROOT_PASSWORD" WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
ALTER USER 'root'@'localhost' IDENTIFIED BY '';
EOF

  if [ "$MYSQL_DATABASE" != "" ]; then
    echo "[i] Creating database: $MYSQL_DATABASE"
    echo "CREATE DATABASE IF NOT EXISTS \`$MYSQL_DATABASE\` CHARACTER SET utf8 COLLATE utf8_general_ci;" >> $tfile

    if [ "$MYSQL_USER" != "" ]; then
      echo "[i] Creating user: $MYSQL_USER with password $MYSQL_PASSWORD"
      echo "GRANT ALL ON \`$MYSQL_DATABASE\`.* to '$MYSQL_USER'@'localhost' IDENTIFIED BY '$MYSQL_PASSWORD';" >> $tfile
    fi
  fi

  /usr/bin/mysqld --verbose --init-file=$tfile &
fi