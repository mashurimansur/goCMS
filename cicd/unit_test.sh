set -ex
echo "Proceeding with Unit tests..."
# go install gotest.tools/gotestsum@v1.12.3

gotestsum --junitfile report.xml --format testname --format dots -- -cover -covermode=count -coverprofile=coverage.out.temp $(go list ./... | grep -v ./integration_test) && cat coverage.out.temp | grep -v -e "_mock.go" > coverage.out && go tool cover -func coverage.out | grep 'total' | sed -e 's/\t\+/ /g;s/%//'| awk '{print $3}' && rm coverage.out.temp coverage.out report.xml