.PHONY: help build fmt
.DEFAULT_GOAL := help

# this is a magic code to output help message at default
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build dist binaries.
	rm -rf dist
	gox \
	  -os="linux darwin windows" \
	  -arch="amd64" \
	  -ldflags=" -w \
	    -X main.CommitHash=`git log --pretty=format:%H -n 1`" \
	  -output "dist/buildsh_{{.OS}}_{{.Arch}}" \
	  .
	cd dist && \
	  mv buildsh_darwin_amd64 buildsh && zip buildsh_darwin_amd64.zip buildsh && rm buildsh && \
	  mv buildsh_linux_amd64 buildsh && zip buildsh_linux_amd64.zip buildsh && rm buildsh && \
	  mv buildsh_windows_amd64.exe buildsh.exe && zip buildsh_windows_amd64.zip buildsh.exe && rm buildsh.exe

dev:  ## Build dev binaru
	go build \
	  -ldflags="-w -X main.CommitHash=`git log --pretty=format:%H -n 1`" \
	  -o="buildsh" .

fmt: ## go fmt
	go fmt $$(go list ./... | grep -v vendor)
	
