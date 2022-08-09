help: ## Display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-16s\033[0m - %s\n", $$1, $$2}'
	@ echo

param-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Param \"$*\" is missing, use: make $(MAKECMDGOALS) $*=<value>"; \
		exit 1; \
	fi

export version?=latest
export COMPOSE_PROJECT_NAME=cadicallegari

build: param-version ## Build the docker image
	docker-compose build

push: param-version ## Push the docker image to our private docker registry
	docker-compose push

# To be used internally or during the development

dev-build: ## Build the dev docker image
	target=builder $(MAKE) build

create-network: ## create network required to run the docker environment
	@docker network create cadicallegari_network 2> /dev/null || exit 0

up: ## Run the service on docker-compose locally
	@docker-compose up -d

down: ## Stop the service on docker-compose locally
	@docker-compose down

export GIT_TAG ?= $(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD)
export GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

MAKE := $(MAKE) -f $(MAKEFILE_LIST)
LDFLAGS = -s -w -X "main.tag=$(GIT_TAG)" -X "main.gitCommit=$(GIT_COMMIT)" -X "main.buildTime=$(shell date -u +%FT%T%z)"

go-generate:
	@go generate ./...

testflag?=-race -cover $(flag)
test: ## Run unit tests, set testcase=<testcase> and flag=-v if you need them
	go test -failfast ./... $(testflag) $(if $(testcase),-run "$(testcase)")

integration-test: override flag+=-tags 'integration' ## Run integration tests, set testcase=<testcase> and flag=-v if you need them
integration-test:
	go test -failfast ./... $(testflag) $(if $(testcase),-run "$(testcase)")

external-test: ;

test-coverage: override flag+=-coverprofile coverage.out ## Show test coverage
test-coverage: coverage.out
	go tool cover -html=coverage.out

coverage.out:
	testflag="$(flag)" $(MAKE) test

go-build: ## Build the binaries
	go build -v -ldflags "$(LDFLAGS)" ./cmd/user

go-install: ## Build the binaries statically and install it
	CGO_ENABLED=0 go install -v -ldflags "$(LDFLAGS)" -a -installsuffix cgo ./cmd/user

run: go-build ## Build and run the app locally
	@./user

cmd?=/bin/sh
exec:
	docker-compose exec user $(cmd)
