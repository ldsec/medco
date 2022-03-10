#!/bin/sh

VERSION=$(git describe --tags 2> /dev/null)
if [ $? -eq 0 ]; then
  echo "$VERSION"
else
  echo "v0.0.0-dev-$(git describe --tags --always)"
fi
