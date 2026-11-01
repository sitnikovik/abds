# Run all tests in the project.
.PHONY: test
test: unit-test integration-test

# Run unit tests.
.PHONY: unit-test
unit-test:
	@echo 🧪 Running unit tests...
	@mkdir -p tmp
	@pkgs=$$(go list ./... 2>/dev/null | grep -v 'internal/test/integration' || true); \
	if [ -n "$$pkgs" ]; then \
		go test -race -count=1 \
		-covermode=atomic \
		-coverprofile=tmp/coverage_unit.out \
		$$pkgs; \
	else \
		echo "⚪ Skipping unit tests: no Go packages found."; \
		echo "mode: atomic" > tmp/coverage_unit.out; \
	fi

# Run integration tests.
.PHONY: integration-test
integration-test:
	@echo 🧪 Running integration tests...
	@mkdir -p tmp
	@set -e; \
	intpkgs=$$(go list -tags=integration ./internal/test/integration/... 2>/dev/null || true); \
	if [ -z "$$intpkgs" ]; then \
		echo "⚪ Skipping integration tests: no integration packages found."; \
		echo "mode: atomic" > tmp/coverage_integration.out; \
		exit 0; \
	fi; \
	coverpkgs=$$(go list ./... 2>/dev/null \
		| grep -v '/internal/test/' \
		| tr '\n' ',' | sed 's/,$$//' || true); \
 	go test -tags=integration \
		-race -count=1 -v \
		-covermode=atomic \
		-coverprofile=tmp/coverage_integration.out \
		-coverpkg=$$coverpkgs \
		$$intpkgs
