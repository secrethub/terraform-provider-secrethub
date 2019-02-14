TEST?=$$(go list ./...)

build:
	go build

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
