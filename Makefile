FIRSTGOPATH = $(firstword $(subst :, ,$(GOPATH)))

gosources = $(shell find . -path "./vendor/*" -prune -o -type f -name "*.go" -print)

all: jasons-game-land-service

jasons-game-land-service: $(gosources) go.mod go.sum
	go build

lint: $(FIRSTGOPATH)/bin/golangci-lint
	$(FIRSTGOPATH)/bin/golangci-lint run --build-tags integration

$(FIRSTGOPATH)/bin/golangci-lint:
	./scripts/download-golangci-lint.sh

test: $(gosources) go.mod go.sum
	go test ./... -tags=integration

install: $(gosources) go.mod go.sum
	go install -a -gcflags=-trimpath=$(CURDIR) -asmflags=-trimpath=$(CURDIR)

clean:
	go clean
	rm -rf vendor

.PHONY: all test clean install lint
