TEST?=$$(go list ./...)

VERSION=`git describe --always`
BUILD_FLAGS=-ldflags "-X "github.com/secrethub/terraform-provider-secrethub/secrethub.version=${VERSION}""

build:
	go build ${BUILD_FLAGS}

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

install:
	go build ${BUILD_FLAGS} -o ~/.terraform.d/plugins/terraform-provider-secrethub

GOLANGCI_VERSION=v1.27.0
lint:
	@docker run --rm -t --user $$(id -u):$$(id -g) -v $$(go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -v $$(go env GOPATH)/pkg:/go/pkg -v ${PWD}:/app -w /app golangci/golangci-lint:${GOLANGCI_VERSION}-alpine golangci-lint run ./...


release-dev:
	env=dev make release

release-prd:
	env=prd make release

PGP_FINGERPRINT=`secrethub read secrethub/terraform-provider/pgp/fingerprint.pgp`

release:
	secrethub read secrethub/terraform-provider/pgp/private.pgp | gpg --import --batch
	PGP_FINGERPRINT=${PGP_FINGERPRINT} goreleaser release --rm-dist $(if $(filter ${env},prd),,--snapshot)
	gpg --yes --batch --delete-secret-key ${PGP_FINGERPRINT}
	gpg --yes --batch --delete-key ${PGP_FINGERPRINT}
