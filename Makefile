.PHONY: test

test:
	go list ./... | grep -v /docs | xargs -n1 go test -cover