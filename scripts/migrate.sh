#!/bin/bash

set -e

CONFIG_FILE="./config/config.json"

cd $PWD

echo $PWD

if [ ! -f $CONFIG_FILE ]; then
  echo "config.json not exist. "
  exit 1
fi

if [[ -z "${MIGRATE_CONNECTION}" ]]; then
  # get database configuration fron config.json
  host=$(cat config/config.json | ./scripts/bash-json-parser.sh | sed -n 's/^database.host=\(.*\)/\1/p')
  port=$(cat config/config.json | ./scripts/bash-json-parser.sh | sed -n 's/^database.port=\(.*\)/\1/p')
  user=$(cat config/config.json | ./scripts/bash-json-parser.sh | sed -n 's/^database.user=\(.*\)/\1/p')
  password=$(cat config/config.json | ./scripts/bash-json-parser.sh | sed -n 's/^database.password=\(.*\)/\1/p')
  database=$(cat config/config.json | ./scripts/bash-json-parser.sh | sed -n 's/^database.dbName=\(.*\)/\1/p')
  connection="postgres://$user:$password@$host:$port/$database?sslmode=disable"
else
  connection="${MIGRATE_CONNECTION}"
fi

./migrate -database $connection -path migrations $@