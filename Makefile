build+darwin+arm64:
	mkdir -p _bin
	GOOS=darwin GOARCH=arm64 go build -o _bin/coda-darwin-arm64 main.go
	chmod 755 _bin/coda-darwin-arm64

build: build+darwin+arm64

install+local: build+darwin+arm64
	sudo cp _bin/coda-darwin-arm64 /usr/local/bin/coda

test:
	go test ./...

dev+json:
	go run main.go jj test.coda.json

dev+yaml:
	go run main.go yy test.coda.yaml

air: 
	air

coverage:
	rm -rf _cover/*
	mkdir -p _cover
	go test -v -coverprofile _cover/coda.cover.out .
	go tool cover -html _cover/coda.cover.out -o _cover/coda.cover.html
	open _cover/coda.cover.html