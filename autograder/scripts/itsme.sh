#!/usr/bin/env sh

# Script to replace placeholders in specified files with provided GitHub nickname and repository name.
# Usage: ./itsme.sh <github_nickname> <github_repo_name>
# If no arguments provided, script will prompt for input interactively.

# Imports
CURDIR=$(dirname "$0")
. "$CURDIR/prompt_var.sh"

# Get arguments or prompt for them.
GITHUB_NICKNAME=$1
GITHUB_REPO=$2
if [ -z "$GITHUB_NICKNAME" ] || [ -z "$GITHUB_REPO" ]; then
  DEFAULT_NICKNAME="${GITHUB_NICKNAME:-$(git config --global user.github 2>/dev/null)}"
  DEFAULT_REPO="${GITHUB_REPO:-$(basename "$(pwd)")}"
  if [ -z "$GITHUB_NICKNAME" ]; then
    prompt_var GITHUB_NICKNAME "Enter GitHub nickname" "$DEFAULT_NICKNAME"
  fi
  if [ -z "$GITHUB_REPO" ]; then
    prompt_var GITHUB_REPO "Enter repository name" "$DEFAULT_REPO"
  fi
fi

# Replace placeholders in files.
FILES="README.md CONTRIBUTING.md"
for FILE in $FILES; do
  if [ ! -f "$FILE" ]; then
    echo "❌ File '$FILE' does not exist."
    exit 1
  fi
  # find {github_nickname} or {github_repo} placeholders in file and replace them with args
    if grep -q "{github_nickname}" "$FILE"; then
        sed -i.bak "s/{github_nickname}/$GITHUB_NICKNAME/g" "$FILE"
        rm "$FILE.bak"
    fi
    if grep -q "{github_repo_name}" "$FILE"; then
        sed -i.bak "s/{github_repo_name}/$GITHUB_REPO/g" "$FILE"
        rm "$FILE.bak"
    fi
done
