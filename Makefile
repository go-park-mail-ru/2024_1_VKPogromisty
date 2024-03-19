.PHONY: test
.PHONY: coverage

test:
	go test ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "docs" > cover.out

coverage:
	go tool cover -func cover.out

docker-build:
	docker-compose build

docker-run:
	docker-compose up

docker-migrate:
	cd ./internal/repository/postgres/migrations && tern migrate
	