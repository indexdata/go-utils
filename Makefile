GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
COVERAGE=coverage.out

.PHONY: check fmt fmt-check vet clean

check:
	$(GO) test -v -cover -coverpkg=./... -coverprofile=$(COVERAGE) ./...

fmt:
	$(GOFMT) -w $(GOFILES)

fmt-check:
	$(GOFMT) -d $(GOFILES)

vet:
	$(GO) vet $(PACKAGES)

clean:
	rm -f $(COVERAGE)
