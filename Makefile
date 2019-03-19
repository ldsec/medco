# generate/update go server based on swagger specifications
swagger-gen:
	swagger validate ./swagger/swagger.yml
	swagger generate server \
		--target=./swagger/ \
		--spec=./swagger/swagger.yml \
		--name=medco-connector
