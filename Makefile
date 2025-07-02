test: 
	go test -v ./...

dependencies:
	./scripts/dependencies.sh

coverage: dependencies
	./scripts/coverage.sh
