#!/bin/bash

db_container=crawler-db
broker_container=crawler-mq

# start database

# Clean previous container
docker ps --filter name=$db_container -q | xargs -I {} docker stop {}

docker run --name $db_container --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres

# Wait for postgres container
sleep 2

docker exec $db_container psql -U postgres -c "create database crawler"

sleep 2

docker exec $db_container psql -U postgres -d crawler -c "
create table if not exists links (
    parent text,
    url text primary key,
    crawled_at timestamp with time zone default clock_timestamp()
)
"

# start broker
docker ps --filter name=$broker_container -q | xargs -I {} docker stop {}

docker run --name $broker_container --rm -d -p 5672:5672 rabbitmq
