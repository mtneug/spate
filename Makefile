# Copyright (c) 2016 Matthias Neugebauer <mtneug@mailbox.org>
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
BUILD_DATE=$(shell date -u +"%Y-%m-%d %T %Z")

PKG=$(shell cat .godir)
PKG_INTEGRATION=${PKG}/integration
PKGS=$(shell go list ./... | grep -v /vendor/)

GO_LDFLAGS=-ldflags " \
	-s \
	-X '$(PKG)/version.gitCommit=$(GIT_COMMIT)' \
	-X '$(PKG)/version.gitTreeState=$(GIT_TREE_STATE)' \
	-X '$(PKG)/version.buildDate=$(BUILD_DATE)'"
GO_BUILD_ARGS=-v $(GO_LDFLAGS)

GOMETALINTER_COMMON_ARGS=\
	--sort=path \
	--vendor \
	--tests \
	--vendored-linters \
	--disable-all \
	--enable=gofmt \
	--enable=vet \
	--enable=vetshadow \
	--enable=golint \
	--enable=ineffassign \
	--enable=goconst \
	--enable=goimports \
	--enable=staticcheck \
	--enable=unused \
	--enable=misspell \
	--enable=lll \
	--line-length=120

all: lint build test integration
ci: lint-full build coverage coverage-integration

build:
	@echo "🌊  $@"
	@go build $(GO_BUILD_ARGS) -o bin/spate $(PKG)

install:
	@echo "🌊  $@"
	@go install $(GO_BUILD_ARGS) $(PKG)

run: build
	@echo "🌊  $@"
	@bin/spate \
		--log-level debug \
		--controller-period 1s \
		--default-autoscaler-period 5s \
		--default-observer-period 5s \
		--default-cooldown-scaled_down 2s \
		--default-cooldown-scaled_up 2s \
		--default-cooldown-service_added 2s \
		--default-cooldown-service_updated 2s

clean:
	@echo "🌊  $@"
	@rm -f bin

lint:
	@echo "🌊  $@"
	@test -z "$$(gometalinter --deadline=5s ${GOMETALINTER_COMMON_ARGS} ./... | tee /dev/stderr)"

lint-full:
	@echo "🌊  $@"
	@test -z "$$(gometalinter --deadline=5m ${GOMETALINTER_COMMON_ARGS} \
			--enable=deadcode \
			--enable=varcheck \
			--enable=structcheck \
			--enable=errcheck \
			--enable=unconvert \
			./... | \
		tee /dev/stderr)"

test:
	@echo "🌊  $@"
	@go test -parallel 8 -race $(filter-out ${PKG_INTEGRATION},${PKGS})

integration:
	@echo "🌊  $@"
	@go test -parallel 8 -race ${PKG_INTEGRATION}

coverage:
	@echo "🌊  $@"
	@status=0; \
	for pkg in $(filter-out ${PKG_INTEGRATION},${PKGS}); do \
		go test -race -coverprofile="../../../$$pkg/coverage.txt" -covermode=atomic $$pkg; \
		true $$((status=status+$$?)); \
	done; \
	exit $$status

coverage-integration:
	@echo "🌊  $@"
	@go test -race -coverprofile="../../../${PKG_INTEGRATION}/coverage.txt" -covermode=atomic ${PKG_INTEGRATION}

ci-docker-image-release:
	@echo "🌊  $@"
	@git clone --depth 1 git@github.com:mtneug/spate-docker.git ../spate-docker

	@# Commit binary
	@echo "Commit binary"
	@cp bin/spate ../spate-docker/spate
	@../spate-docker/update-image.sh "${TRAVIS_TAG}" "${TRAVIS_COMMIT}"

	@cd ../spate-docker && git add -A
	@cd ../spate-docker && git commit -m "Release ${TRAVIS_TAG} - ${TRAVIS_COMMIT}"
	@cd ../spate-docker && git tag -f "${TRAVIS_TAG}"

	@# Update README.md
	@echo "Update README"
	@../spate-docker/update-readme.sh

	@cd ../spate-docker && git add README.md
	@cd ../spate-docker && git commit -m "Update README.md"

	# Push
	@cd ../spate-docker && git push -f
	@cd ../spate-docker && git push -f --tags

.PHONY: all ci build install clean lint lint-full test integration coverage coverage-integration ci-docker-image-release
