#!/usr/bin/env bash
#
# This script calculates coverage percentage from a coverage file.
# Usage: calculate_coverage.sh <coverage_file> <output_file> <label> [include_pattern]

set -e

coverage_file=$1
output_file=$2
label=$3
include_pattern=$4

if [ -z "$coverage_file" ] || [ -z "$output_file" ] || [ -z "$label" ]; then
	echo "Usage: $0 <coverage_file> <output_file> <label> [include_pattern]" >&2
	exit 1
fi

# Filter out main packages and calculate coverage percentage.
filtered_file="${coverage_file}.filtered"
head -n1 "${coverage_file}" > "${filtered_file}"
if [ -n "$include_pattern" ]; then
	grep -v '/cmd/' "${coverage_file}" | grep "$include_pattern" | tail -n +2 >> "${filtered_file}"
else
	grep -v '/cmd/' "${coverage_file}" | grep -v '/internal/test/' | tail -n +2 >> "${filtered_file}"
fi

percent=$(go tool cover -func="${filtered_file}" 2>/dev/null \
	| tail -n1 \
	| awk '{print $NF}')

rm -f "${filtered_file}"

# Remove the % sign from the end of percent
percent_no_pct="${percent%\%}"
printf '%s' "$percent_no_pct" > "$output_file"
echo "${label}: $percent"
