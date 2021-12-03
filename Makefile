BUILD_COMMIT := $(shell git describe --tags --always --dirty --match=v*)
BUILD_DATE := $(shell date -u +%b-%d-%Y,%T-UTC)
BUILD_SEMVER := $(shell cat .SEMVER)

.PHONY: all clean release dirty-check test help

# target: all - default target, will trigger build
all: test

# target: clean - removes all build and test artifacts
clean:
	@-rm -rf results

# target: test - runs tests and generates coverage reports
test:
	mkdir -p results
	go test ./... -race -cover -coverprofile=results/tc.out
	go tool cover -html=results/tc.out -o results/coverage.html

# target: release - will clean, test, and finally creates a git tag for the version
release: dirty-check clean test
	git tag v$(BUILD_SEMVER) $(BUILD_COMMIT)
	git push origin v$(BUILD_SEMVER)

# target: dirty-check - will check if repo is dirty
dirty-check:
ifneq (, $(findstring dirty, $(BUILD_COMMIT)))
	@echo "you're dirty check your repo status before releasing"
	false
endif

# target: help - displays help
help:
	@egrep "^#.?target:" Makefile