TEST?=$$(go list ./...)

VERSION=`git describe --always`
BUILD_FLAGS=-ldflags "-X "github.com/secrethub/terraform-provider-secrethub/secrethub.version=${VERSION}""

build:
	go build ${BUILD_FLAGS}

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

install:
	go build ${BUILD_FLAGS} -o ~/.terraform.d/plugins/terraform-provider-secrethub

GOLANGCI_VERSION=v1.23.8

lint:
	docker run -t --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.23.8-alpine golangci-lint run --timeout 2m0s ./...
