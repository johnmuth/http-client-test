#!/bin/bash -e

function lintAppCode() {
    gometalinter \
        --vendor \
        --exclude='error return value not checked.*(Close|Log|Print).*\(errcheck\)$' \
        --exclude='.*_test\.go:.*error return value not checked.*\(errcheck\)$' \
        --exclude='duplicate of.*_test.go.*\(dupl\)$' \
        --disable=aligncheck \
        --disable=golint \
        --disable=gotype \
        --disable=structcheck \
        --disable=varcheck \
        --disable=unconvert \
        --disable=aligncheck \
        --disable=dupl \
        --disable=goconst  \
        --disable=gosimple  \
        --disable=staticcheck \
        --cyclo-over=20 \
        --tests \
        --deadline=30s
}

function fmtAppCode() {
    (
        cd "$PROJECT_BASE_DIR"
        go fmt $(go list ./... | grep -v /vendor/)
    )
}

function runTests() {
    # We need to do a bit of fiddling to generate coverage data from multiple packages and merge them.
    (
        if [ ! -d coverage ]; then
            mkdir coverage
        fi
        chmod 777 coverage
        cd "$PROJECT_BASE_DIR"
        echo "mode: count" > coverage/coverage-all.out
        for pkg in $(go list ./... | grep -v /vendor/)
        do
          echo "pkg=$pkg"
          richgo test -v -coverprofile=coverage/coverage.out -covermode=count $pkg
          if [ -f coverage/coverage.out ]; then
            tail -n +2 coverage/coverage.out >> coverage/coverage-all.out
          fi
        done
        richgo test  -v -tags=integration -coverprofile=coverage/coverage.out -covermode=count
        tail -n +2 coverage/coverage.out >> coverage/coverage-all.out
        go tool cover -html=coverage/coverage-all.out -o coverage/coverage.html
    )
}

function buildAppCode() {
    rm -f app
    go build -v ./...
}