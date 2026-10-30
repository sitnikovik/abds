# Lint the codebase.
.PHONY: lint
lint:
	@pkgs=$$(go list ./... 2>/dev/null || true); \
	if [ -z "$$pkgs" ]; then \
		echo "⚪ Skipping lint: no Go packages found."; \
	else \
		golangci-lint run ./...; \
	fi
