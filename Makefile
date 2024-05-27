PROJECT_NAME := nextride-shortcut
PKG := "$(PROJECT_NAME)"
CMD := server.go
GO_FILES := $(shell find . -type f -name '*.go' | grep -v _test.go)
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PROJECT_NAME := $(PROJECT_NAME)-$(OS)-$(ARCH)
BINDIR := bin

.SUFFIXES:
.PHONY: all coverage lint test race x2unit xunit clean download locallint printexe

all: $(BINDIR)/$(PROJECT_NAME) $(BINDIR)/$(PROJECT_NAME).sha256

ifeq ($(GITHUB_ACTIONS),)
coverage: html
lint: locallint
else
coverage: coverage.out
lint: cilint.txt
endif

ifdef GITHUB_HEAD_REF
LDFLAGS := -ldflags="-X main.Version=$(GITHUB_HEAD_REF)"
endif

locallint:
	@golangci-lint run

cilint.txt: $(GO_FILES)
	@golangci-lint run --out-format=line-number --new-from-rev=master --issues-exit-code=0 > $@

test:
	@go test -short ./...

race:
	@go test -v -count=1 -race ./...

coverage.out: $(GO_FILES)
	@go test -v -race -covermode=atomic -coverprofile=$@ ./...

html: coverage.out
	@go tool cover --html=$<

x2unit:
	go get github.com/tebeka/go2xunit

tests.out:
	@go test -v -race ./... > $@

xunit: x2unit tests.out
	go2xunit -fail -input tests.out -output tests.xml
	@rm -f tests.out

run:
	@go run $(CMD) --verbosity debug

clean:
	@rm -rf tests.xml tests.out coverage.out *.sha256 main cilint.txt $(BINDIR)/ $(PKG) $(PKG)-*
	@docker rm --force $(PKG) &> /dev/null || true
	@go mod tidy

$(BINDIR)/$(PROJECT_NAME): $(GO_FILES)
	@go build -v -o $@ $(LDFLAGS) $(CMD)

%.sha256:
	@openssl dgst -sha256 -hex $* | cut -f2 -d' ' > $@

download:
	@go mod download

printexe:
	@echo $(PROJECT_NAME)

