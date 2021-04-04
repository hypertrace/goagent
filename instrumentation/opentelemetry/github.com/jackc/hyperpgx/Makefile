.PHONY: test-integration
test-integration:
	docker run --rm --name test_postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres
	@sleep 1
	go test --tags=integration ./... || { docker rm -f test_postgres ; exit 1 ; }
	@docker rm -f test_postgres