#!/usr/bin/env bash
PGPASSWORD="i2b2" psql -U i2b2 -d i2b2medcosrv0 -h localhost -p 5432 -a -f sqlAnnotations_all.sql
PGPASSWORD="i2b2" psql -U i2b2 -d i2b2medcosrv1 -h localhost -p 5432 -a -f sqlAnnotations_all.sql
PGPASSWORD="i2b2" psql -U i2b2 -d i2b2medcosrv2 -h localhost -p 5432 -a -f sqlAnnotations_all.sql