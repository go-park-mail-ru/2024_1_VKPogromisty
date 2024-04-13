.PHONY: test
.PHONY: coverage
.PHONY: mocks

MOCKS_DESTINATION=mocks
mocks: 
	@echo "Generating mocks...";
	@rm -rf $(MOCKS_DESTINATION);
	@for file in $(shell find usecase -type f -name '*.go' ! -name '*_test.go'); do \
		mkdir -p $(MOCKS_DESTINATION)/`dirname $$file` && mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; \
	done

test:
	go test ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "docs" | grep -v "mocks" > cover.out

coverage:
	go tool cover -func cover.out

docker-build:
	docker-compose build --no-cache

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-migrate:
	cd ./db/migrations && tern migrate

make go-run:
	cd ./cmd && go build && ./cmd
