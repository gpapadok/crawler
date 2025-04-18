#!/bin/bash

db_container=crawler-db
broker_container=crawler-mq
seed_container=crawler-seed

ROOT_DIR=$(dirname $(dirname $0))

source $ROOT_DIR/.env

# create docker network
docker network create --driver bridge crawler-net

# start database

# Clean previous container
docker ps --filter name=$db_container -q | xargs -I {} docker stop {}

docker run --rm -d \
       --network crawler-net \
       -p 5432:5432 \
       -e POSTGRES_PASSWORD=$DB_PASSWORD \
       -v $(dirname $(cd $ROOT_DIR; pwd))/docker/pgdata:/var/lib/postgresql/data \
       --name $db_container postgres

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

docker run --network crawler-net --name $broker_container --rm -d -p 5672:5672 rabbitmq

# Wait for rabbit-mq container
sleep 2

# seed database and broker
docker run --network crawler-net --rm $seed_container
