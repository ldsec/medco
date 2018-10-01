#!/usr/bin/env bash

docker build .
docker run --volume=build:/opt/build-dir