#!/usr/bin/env bash
./medco-loader -debug 2 v1 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entry 0 --sen ../data/i2b2/sensitive.txt -f ../data/i2b2/files.toml \
--dbHost localhost  --dbPort 5432 --dbName i2b2medcosrv0 --dbUser i2b2 --dbPassword i2b2

./medco-loader -debug 2 v1 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entry 1 --sen ../data/i2b2/sensitive.txt -f ../data/i2b2/files.toml \
--dbHost localhost  --dbPort 5432 --dbName i2b2medcosrv1 --dbUser i2b2 --dbPassword i2b2

./medco-loader -debug 2 v1 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entry 2 --sen ../data/i2b2/sensitive.txt -f ../data/i2b2/files.toml \
--dbHost localhost  --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2