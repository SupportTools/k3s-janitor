#!/bin/sh

echo "Setting docker environment"
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD harbor.support.tools

docker pull supporttools/k3s-janitor:latest
echo "Building..."
if ! docker build -t supporttools/k3s-janitor:${DRONE_BUILD_NUMBER} --cache-from supporttools/k3s-janitor:latest -f Dockerfile .
then
  echo "Docker build failed"
  exit 127
fi
if ! docker push supporttools/k3s-janitor:${DRONE_BUILD_NUMBER}
then
  echo "Docker push failed"
  exit 126
fi
echo "Tagging to latest and pushing..."
if ! docker tag supporttools/k3s-janitor:${DRONE_BUILD_NUMBER} supporttools/k3s-janitor:latest
then
  echo "Docker tag failed"
  exit 123
fi
if ! docker push supporttools/k3s-janitor:latest
then
  echo "Docker push failed"
  exit 122
fi