test-all:
	go test --tags=rabbitmq_client_test,eventstoredb_test -v ./...

test-eventstoredb:
	go test --tags=eventstoredb_test -v ./eventstoredb