# Check if the total test coverage meets the defined threshold.
.PHONY: coverage
coverage:
	@COVERAGE_THRESHOLD=$(COVERAGE_THRESHOLD) sh scripts/coverage/check.sh

# Generate coverage report from unit tests.
# Outputs total coverage report to tmp/coverage_total.out,
# covered lines to tmp/coverage_*.out
# and uncovered lines to tmp/uncovered.out.
.coverage-report:
	@sh scripts/coverage/report.sh
