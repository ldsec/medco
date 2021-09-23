export MEDCO_VERSION := $(shell make --no-print-directory -C ../ medco_version)

.PHONY: pull up stop down
pull:
	docker-compose -f docker-compose.yml pull

up:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down