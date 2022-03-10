#!/bin/bash

if [[ -n "$1" ]]; then
  # if provided, use first argument for version
  MEDCO_VERSION=$1

elif [[ $GITHUB_REF == refs/tags/* ]]; then
  # if in the CI and under a tag, override
  MEDCO_VERSION=${GITHUB_REF#refs/tags/}

else
  # get version from git describe
  MEDCO_VERSION=$(git describe --tags 2> /dev/null)

  # if failed because of no available tag, only use commit
  if [[ $? -ne 0 ]]; then
    MEDCO_VERSION="v0.0.0-dev-$(git describe --tags --always)"
  fi
fi
export MEDCO_VERSION

export MEDCO_DOCKER_TAG=ghcr.io/chuv-ds/medco:$MEDCO_VERSION
export I2B2_DOCKER_TAG=ghcr.io/chuv-ds/i2b2-medco:$MEDCO_VERSION
export KEYCLOAK_DOCKER_TAG=ghcr.io/chuv-ds/keycloak-medco:$MEDCO_VERSION
export WSTUNNEL_DOCKER_TAG=ghcr.io/chuv-ds/medco-unlynx-wstunnel:$MEDCO_VERSION
export NGINX_DOCKER_TAG=ghcr.io/chuv-ds/nginx-medco:$MEDCO_VERSION
export PGADMIN_TAG=ghcr.io/chuv-ds/pgadmin-medco:$MEDCO_VERSION

export GB_DOCKER_TAG=ghcr.io/ldsec/glowing-bear-medco:v3.0.0
export PG_DOCKER_TAG=postgres:9.6
