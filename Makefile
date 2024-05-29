PROJECT_NAME := nextride-shortcut
PKG := "$(PROJECT_NAME)"
CMD := server.go
GO_FILES := $(shell find . -type f -name '*.go' | grep -v _test.go)
TEST_FILES := $(shell find . -type f -name '*_test.go')
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PROJECT_NAME := $(PROJECT_NAME)-$(OS)-$(ARCH)
BINDIR := bin

.SUFFIXES:
.PHONY: all coverage lint test junit clean download docker

all: $(BINDIR)/$(PROJECT_NAME) $(BINDIR)/$(PROJECT_NAME).sha256

ifdef GITHUB_HEAD_REF
LDFLAGS := -ldflags="-X main.Version=$(GITHUB_HEAD_REF)"
endif

lint:
	@golangci-lint run --out-format=line-number

test:
	@go test -short ./...

coverage.out: $(GO_FILES) $(TEST_FILES)
	@go test -v -race -covermode=atomic -coverprofile=$@ ./... 2>&1 | tee tests.out

coverage: coverage.out
ifeq ($(GITHUB_ACTIONS),)
	@go tool cover --html=$<
else
	@go tool cover --html=$< -o coverage.html
endif

junit:
	@go install github.com/jstemmer/go-junit-report@latest

report.xml: coverage
	@go-junit-report -set-exit-code < tests.out > $@

run:
	@go run $(CMD) --verbosity debug

clean:
	@rm -rf tests.xml *.out report.xml *.sha256 main cilint.txt $(BINDIR)/ $(PKG) $(PKG)-*
	@docker rm --force $(PKG) &> /dev/null || true
	@go mod tidy

$(BINDIR)/$(PROJECT_NAME): $(GO_FILES)
	@go build -v -o $@ $(LDFLAGS) $(CMD)

%.sha256:
	@openssl dgst -sha256 -hex $* | cut -f2 -d' ' > $@

download:
	@go mod download

docker: $(BINDIR)/$(PROJECT_NAME)
	@docker build -t $(PROJECT_NAME):dev --build-arg GOOS=$(OS) --build-arg GOARCH=$(ARCH) .
