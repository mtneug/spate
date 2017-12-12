# Copyright (c) 2016 Matthias Neugebauer
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

PKGS=$(shell go list ./...)

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

all: lint test
ci: lint-full coverage

lint:
	@echo "$@"
	@test -z "$$(gometalinter --deadline=5s ${GOMETALINTER_COMMON_ARGS} ./... | tee /dev/stderr)"

lint-full:
	@echo "$@"
	@test -z "$$(gometalinter --deadline=5m ${GOMETALINTER_COMMON_ARGS} \
			--enable=deadcode \
			--enable=varcheck \
			--enable=structcheck \
			--enable=errcheck \
			--enable=unconvert \
			./... | \
		tee /dev/stderr)"

test:
	@echo "$@"
	@go test -parallel 8 -race ${PKGS}

coverage:
	@echo "$@"
	@status=0; \
	for pkg in ${PKGS}; do \
		go test -race -coverprofile="../../../$$pkg/coverage.txt" -covermode=atomic $$pkg; \
		true $$((status=status+$$?)); \
	done; \
	exit $$status

.PHONY: all ci lint lint-full test coverage
