.PHONY: test
.PHONY: coverage
.PHONY: mocks

MOCKS_DESTINATION=mocks
mocks: 
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $(shell find usecase -type f -name '*.go' ! -name '*_test.go' ! -name '*_easyjson.go'); do \
		mkdir -p $(MOCKS_DESTINATION)/`dirname $$file` && mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; \
	done
	@for file in $(shell find internal/rest/chat -type f -name '*.go' ! -name '*_test.go' ! -name '*_easyjson.go'); do \
		mkdir -p $(MOCKS_DESTINATION)/`dirname $$file | sed 's/internal\///'` && mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/`echo $$file | sed 's/internal\///'`; \
	done
	@mkdir -p mocks/grpc
	@find internal/grpc -name '*.proto' -print0 | while IFS= read -r -d '' file; do \
		service=$$(dirname "$$file"); \
        service=$${service%/proto}; \
		service_name=$$(basename "$$service"); \
		source_name=$$service_name"_grpc.pb.go"; \
		mockgen -source=$$service"/proto/"$$source_name -destination="mocks/grpc/"$$service_name"_grpc/"$$service_name"_mock.go" -package=$$service_name"_grpc"; \
	done
	@echo "Mocks generated."

EASYJSON_DIRS=./usecase ./domain ./errors ./internal/rest ./pkg/json

.PHONY: clean_easyjson

clean_easyjson:
	@echo "Cleaning up easyjson files..."
	@for dir in $(EASYJSON_DIRS); do \
		for file in $$(find $$dir -type f -name '*_easyjson.go'); do \
			rm $$file; \
		done \
	done
	@echo "Cleanup completed."

.PHONY: easyjson

easyjson:
	make clean_easyjson
	@echo "Generating easyjson..."
	@for dir in $(EASYJSON_DIRS); do \
		for file in $$(find $$dir -type f -name '*.go' ! -name '*_test.go'); do \
			easyjson $$file; \
		done \
	done
	@echo "easyjson generation completed."

test:
	go test ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "docs" | grep -v "mocks" | grep -v "proto" | grep -v "easyjson" > cover.out

coverage:
	go tool cover -func cover.out

swaggen:
	swag init -g cmd/app/main.go

post-build:
	docker build -t socio/post-service -f cmd/post/Dockerfile . --no-cache

user-build:
	docker build -t socio/user-service -f cmd/user/Dockerfile . --no-cache

auth-build:
	docker build -t socio/auth-service -f cmd/auth/Dockerfile . --no-cache

public-group-build:
	docker build -t socio/public-group-service -f cmd/public_group/Dockerfile . --no-cache

app-build:
	docker build -t socio/app-service -f cmd/app/Dockerfile . --no-cache

docker-build:
	make user-build
	make post-build
	make auth-build
	make public-group-build
	make app-build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-migrate:
	cd ./db/migrations && tern migrate

make go-run:
	cd ./cmd && go build && ./cmd
