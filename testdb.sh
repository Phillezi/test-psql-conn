#!/bin/sh

docker run --rm -it -p 5432:5432 -e POSTGRES_USER=myuser -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=mydb -v $(pwd)/init:/docker-entrypoint-initdb.d postgres:15-alpine 
