# To build the project, run `make build`
build:
	@go build

# To run the UTs, run `make test`
test:
	@go test -v ./...