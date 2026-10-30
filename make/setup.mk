.PHONY: setup
setup: itsme githooks

# Install Git hooks for automatic testing and linting.
.PHONY: githooks
githooks:
	@git config core.hooksPath .githooks
	@chmod +x .githooks/pre-commit .githooks/pre-push
	@echo "✅ Git hooks installed successfully"

# Replace all template substrings with developer credentials.
itsme:
	@sh ./scripts/itsme.sh