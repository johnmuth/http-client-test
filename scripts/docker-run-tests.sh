#!/bin/bash -e

docker-compose down --remove-orphans
docker-compose build web
docker-compose run web ./scripts/run-tests.sh
