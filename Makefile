COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

include sqlc.mk

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
