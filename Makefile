test-all:
	go test -count=1 -v --tags=rabbitmq_client_test,eventstoredb_test ./...

test-eventstoredb:
	go test -count=1 -v --tags=eventstoredb_test ./eventstoredb