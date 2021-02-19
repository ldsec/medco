#!/usr/bin/env bash
set -Eeuo pipefail

for script in test/test_cli_preloading/*
do
  bash "$script"
done
echo "CLI test 1/2 successful!"

