VERSION ?= $(shell cat ./pkg/version/.version)

build+darwin+arm64:
	mkdir -p _bin
	GOOS=darwin GOARCH=arm64 go build -o _bin/coda-darwin-arm64 main.go
	chmod 755 _bin/coda-darwin-arm64
build+linux+amd64:
	mkdir -p _bin
	GOOS=linux GOARCH=amd64 go build -o _bin/coda-linux-amd64 main.go
	chmod 755 _bin/coda-linux-amd64

build: build+darwin+arm64
build+linux: build+linux+amd64

install+local: build+darwin+arm64
	sudo cp _bin/coda-darwin-arm64 ~/code/go/bin/coda

test:
	go test ./...

dev+server:
	go run -race main.go server

dev+json:
	go run main.go jj test.coda.json

dev+yaml:
	go run main.go yy test.coda.yaml

build+docker: build+linux+amd64
	docker buildx build --platform linux/amd64 -t ghcr.io/yosev/coda:$(VERSION) -t ghcr.io/yosev/coda:latest -f Dockerfile.local .
	docker push ghcr.io/yosev/coda:$(VERSION) 
	docker push ghcr.io/yosev/coda:latest

air: 
	air

coverage:
	rm -rf _cover/*
	mkdir -p _cover
	go test -v -coverprofile _cover/coda.cover.out .
	go tool cover -html _cover/coda.cover.out -o _cover/coda.cover.html
	open _cover/coda.cover.html