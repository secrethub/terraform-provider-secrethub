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

lint-tools:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION}

lint:
	@golangci-lint run
