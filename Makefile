run: clean test format vet plugins main
	./main

plugins:
	go build -o plugins/ -buildmode=plugin ./plugins_src/*

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
