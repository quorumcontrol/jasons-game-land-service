FIRSTGOPATH = $(firstword $(subst :, ,$(GOPATH)))

gosources = $(shell find . -path "./vendor/*" -prune -o -type f -name "*.go" -print)
generated =messages/external_gen.go messages/external_gen_test.go

all: build

$(generated): messages/external.go $(FIRSTGOPATH)/bin/msgp
	cd messages && go generate

build: $(gosources) $(generated) go.mod go.sum
	go build

lint: $(FIRSTGOPATH)/bin/golangci-lint
	$(FIRSTGOPATH)/bin/golangci-lint run --build-tags integration

$(FIRSTGOPATH)/bin/golangci-lint:
	./scripts/download-golangci-lint.sh

$(FIRSTGOPATH)/bin/msgp:
	go get github.com/tinylib/msgp

test: $(gosources) $(generated) go.mod go.sum
	go test ./... -tags=integration

install: $(gosources) $(generated) go.mod go.sum
	go install -a -gcflags=-trimpath=$(CURDIR) -asmflags=-trimpath=$(CURDIR)

clean:
	go clean
	rm -rf vendor

.PHONY: all build test clean install lint
