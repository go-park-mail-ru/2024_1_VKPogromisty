.PHONY: test

test:
	go test ./... -coverprofile cover.out.tmp | grep -v "docs"
	cat cover.out.tmp > cover.out
