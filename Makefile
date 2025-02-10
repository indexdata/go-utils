GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
COVERAGE=coverage.out

.PHONY: check fmt fmt-check vet clean

check:
	$(GO) test -v -cover -coverpkg=./... -coverprofile=$(COVERAGE) ./...

view-coverage: check
	$(GO) tool cover -html=$(COVERAGE)

fmt:
	$(GOFMT) -w $(GOFILES)

fmt-check:
	$(GOFMT) -d $(GOFILES)

vet:
	$(GO) vet $(PACKAGES)

clean:
	rm -f $(COVERAGE)
