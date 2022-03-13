# Copyright 2018 The Service Manager Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: build test

TEST_PROFILE ?= $(CURDIR)/profile.cov
COVERAGE ?= $(CURDIR)/coverage.html

PLATFORM ?= linux
ARCH     ?= amd64

COVER_OUT			?= $(CURDIR)/cover.out

BUILD_LDFLAGS =

# GO_FLAGS - extra "go build" flags to use - e.g. -v (for verbose)
GO_BUILD = env CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) \
           go build $(GO_FLAGS) -ldflags '-s -w $(BUILD_LDFLAGS)' ./...

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
#TEST_PROFILE=cover.out

build: smcli

goget:
	@go get -v -t -d ./...

smcli:
	$(GO_BUILD)

test: build
	@echo Running tests:
	@go test ./... -p 1 -coverpkg $(shell go list ./... | egrep -v "fakes|test" | paste -sd "," -) -coverprofile=$(TEST_PROFILE)

coverage: test
	@go tool cover -html=$(TEST_PROFILE) -o "$(COVERAGE)"

clean: clean-test clean-coverage

clean-test:
	rm -f $(TEST_PROFILE)

clean-coverage:
	rm -f $(COVERAGE)

clean-vendor:
	rm -rf vendor
	@echo > go.sum

go-deps:
	set GO111MODULE=off
	go get gotest.tools/gotestsum
	go get github.com/t-yuki/gocover-cobertura
	go install github.com/axw/gocov/gocov@latest
	go get github.com/AlekSi/gocov-xml
	go get -u github.com/jstemmer/go-junit-report
	set GO111MODULE=on
	go mod tidy
# Run tests

run-test: go-deps
	@echo COVER_OUT
	rm -rf reports
	mkdir -p reports
	gotestsum --junitfile reports/junit.xml -- -coverprofile=cover.out ./... -mod=mod
	go tool cover -func $(COVER_OUT) | grep total
	find . -name cover.out -execdir sh -c 'gocover-cobertura < cover.out > coverage.xml'  \; ;\
	GO111MODULE=on


