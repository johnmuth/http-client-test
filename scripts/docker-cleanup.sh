#!/bin/bash

# despite of calling `docker-compose build --force-rm` and
# `docker-compose run --rm` docker still produces dangling images and volumes
IMAGES_TO_REMOVE=$(docker images -qf dangling=true)
if [[ -n "$IMAGES_TO_REMOVE" ]]
then
  echo "Removing dangling docker images..."
  docker rmi -f $IMAGES_TO_REMOVE
fi

VOLUMES_TO_REMOVE=$(docker volume ls -qf dangling=true)
if [[ -n "$VOLUMES_TO_REMOVE" ]]
then
  echo "Removing dangling docker volumes..."
  docker volume rm $VOLUMES_TO_REMOVE
fi
