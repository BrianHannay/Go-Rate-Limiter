run: clean test format vet plugins main
	./main

plugins:
	cd plugins && go build -buildmode=plugin ../plugins_src/*

main: plugins main.go
	go build main.go
	chmod +x main

format:
	find . -type f -name '*.go' -exec gofmt -w -e -s -d {} \;

vet:
	go vet ./...

test:
	go test ./test -v

clean:
	rm -f main

.PHONY: run clean test format vet plugins
