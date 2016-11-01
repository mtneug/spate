# Copyright © 2016 Matthias Neugebauer <mtneug@mailbox.org>
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

GIT_COMMIT=$(shell git rev-parse --short HEAD || echo "unknown")
GIT_TREE_STATE=$(shell sh -c 'if test -z "`git status --porcelain 2>/dev/null`"; then echo clean; else echo dirty; fi')
BUILD_DATE=$(shell date -u +"%Y-%m-%d")

GO_PKG=github.com/mtneug/spate
GO_PKG_ALL=$(shell go list ./... | grep -v /vendor/)
GO_LDFLAGS=-ldflags "-X $(GO_PKG)/version.gitCommit=$(GIT_COMMIT) -X $(GO_PKG)/version.gitTreeState=$(GIT_TREE_STATE) -X $(GO_PKG)/version.buildDate=$(BUILD_DATE)"
GO_BUILD_ARGS=-v $(GO_LDFLAGS)

all: build

build:
	@echo "🌊 $@"
	@go build $(GO_BUILD_ARGS) -o bin/spate $(GO_PKG)

install:
	@echo "🌊 $@"
	@go install $(GO_BUILD_ARGS) $(GO_PKG)

test:
	@echo "🌊 $@"
	@go test $(GO_PKG_ALL)

.PHONY: all build install test
