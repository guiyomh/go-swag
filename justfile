
IGNORED_FOLDER := ".ignore"

default:
    @just --list

# Run all tests
test: _create-ignore
    gotestsum --format standard-verbose --jsonfile {{IGNORED_FOLDER}}/report.json --junitfile {{IGNORED_FOLDER}}/report.xml -- -covermode=count -coverprofile={{IGNORED_FOLDER}}/cover.out ./...
    go tool cover -html={{IGNORED_FOLDER}}/cover.out -o {{IGNORED_FOLDER}}/coverage.html

# Cleanup the ignored folders and binaries
clean:
    rm -rf {{IGNORED_FOLDER}}

#Run lint step with golangci-lint
lint:
    golangci-lint run

_create-ignore:
    @if [ ! -d {{IGNORED_FOLDER}} ]; then \
        mkdir -p {{IGNORED_FOLDER}}; \
    fi