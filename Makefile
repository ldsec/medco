validate:
	swagger validate ./swagger.yml

gen: validate
	swagger generate server \
		--target=./ \
		--spec=./swagger.yml \
		--name=medco-connector

.PHONY: install gen validate
