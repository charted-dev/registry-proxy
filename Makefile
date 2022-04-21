# 🪟 registry-proxy: Pluggable Docker Registry proxy for HTTP servers, made for charted-server.
# Copyright 2022 Noelware <team@noelware.org>
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

# Usage: `make help`
.PHONY: help
help: ## Prints the help usage of this Makefile
	@printf "\033[34m────────────────────────────────────────────────────────────────────────────────────────────\033[0m\n"
	@printf ":: Usage ::\n"
	@printf "make <target> [VARIABLE=value]\n"
	@printf "\n:: Targets ::\n"
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z_-]+:.*?##/ { printf "  make \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 25) } ' $(MAKEFILE_LIST)
	@printf "\033[34m────────────────────────────────────────────────────────────────────────────────────────────\033[0m\n"

# Usage: `make test`
.PHONY: test
test: ## Tests the library to see if it's successful.
	@printf "Testing library....\n"
	@go test -v ./...
