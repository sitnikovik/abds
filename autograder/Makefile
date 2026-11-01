.DEFAULT_GOAL = check

-include ../.github/ci.env
include ../make/test.mk
include ../make/lint.mk
include ../make/coverage.mk

# Run all checks required to validate the autograder before merging.
.PHONY: check
check: test coverage lint
