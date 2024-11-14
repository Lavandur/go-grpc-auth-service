#!/bin/bash

function create_db_with_user() {
  local database=$1
  echo "Creating user and database '$database'"
  psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" <<-EOSQL
      CREATE DATABASE "$database" OWNER $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
  echo "Databases for creating: $POSTGRES_MULTIPLE_DATABASES"
  for database in $(echo "$POSTGRES_MULTIPLE_DATABASES" | tr ',' ' ');
  do
    create_db_with_user "$database"
  done
  echo "Databases created"
fi
