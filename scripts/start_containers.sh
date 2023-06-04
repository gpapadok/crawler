#!/bin/bash

db_container=crawler-db
broker_container=crawler-mq

source $(dirname $(dirname $0))/.env

# start database

# Clean previous container
docker ps --filter name=$db_container -q | xargs -I {} docker stop {}

docker run --name $db_container --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=$DB_PASSWORD postgres

# Wait for postgres container
sleep 2

docker exec $db_container psql -U postgres -c "create database crawler"
docker exec $db_container psql -U postgres -c "create database crawler_test"

sleep 2

init_sql=$(cat $(dirname $0)/init.sql)

docker exec $db_container psql -U postgres -d crawler -c "$init_sql"
docker exec $db_container psql -U postgres -d crawler_test -c "$init_sql"

# start broker
docker ps --filter name=$broker_container -q | xargs -I {} docker stop {}

docker run --name $broker_container --rm -d -p 5672:5672 rabbitmq
