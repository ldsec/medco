#!/usr/bin/env bash
bash compileMac.sh
./medcoLoader -debug 2 v1 -g ../data/i2b2/group.toml --entry 0 --sen ../data/i2b2/sensitive.txt -f ../data/i2b2/files.toml -e