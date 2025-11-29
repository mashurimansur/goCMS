#!/bin/bash
set -ex

echo "Proceeding with Unit tests..."
echo "Current PATH: $PATH"
echo "GOPATH: $(go env GOPATH)"

# Fix for Go 1.25+ covdata error
export GOTOOLCHAIN=go1.25.3+auto

# Install gotestsum
go install gotest.tools/gotestsum@v1.12.3

# Configure git for bitbucket (if needed)
git config --global --add url."git@bitbucket.org:".insteadOf "https://bitbucket.org/" || true

# Run tests with coverage
gotestsum \
   --junitfile report.xml --format testname \
   --format dots \
   -- -cover -covermode=count -coverprofile=coverage.out.temp $(go list ./... | grep -v ./integration_test)

# Process coverage results
cat coverage.out.temp | grep -v -e "_mock.go" > coverage.out

# Extract coverage percentage
export CODE_COVERAGE=$(go tool cover -func coverage.out | grep 'total' | sed -e 's/\t\+/ /g;s/%//'| awk '{print $3}')

echo "total: (statements) $CODE_COVERAGE%"
echo code_coverage=$CODE_COVERAGE > code_coverage_results

# Cleanup temp files
rm -f coverage.out.temp coverage.out report.xml