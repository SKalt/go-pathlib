#!/usr/bin/env bash
### USAGE: ./scripts/layout.sh [-h|--help]
### DESCRIPTION: This script is used to expand the directory structure for the project.
###              see https://github.com/golang-standards/project-layout
dirs='
pkg         Public libraries.
internal    Private application and library code.
cmd         Public applications.
scripts     Scripts to perform common operations.
test        Additional external test apps and test data.
web         Web application specific components.
build       Packaging and Continuous Integration.
deployments IaC definitions and scripts for CI/CD and infrastructure.
tools       Support tools and utilities.
'

set -euo pipefail
search=""
declare -a target=()
while [ -n "${1:-}" ]; do
  case "$1" in
    -h|--help) usage && exit 0;;
    *)
      if [ -n "$search" ]; then search="$search $1"; else search="$1"; fi
      shift
      ;;
  esac
done
if [ -n "$search" ]; then
  for t in $(echo "$dirs" | grep "$search" | cut -d' ' -f1); do
    target+=("$t");
  done
else
  for t in $( echo "$dirs" | fzf | cut -d' ' -f1); do
    target+=("$t");
  done
fi
for t in "${target[@]}"; do
  mkdir -p "$t"
  # trim the first column
  msg=$(echo "$dirs" | grep "^$t" | cut -d' ' -f2- | sed 's/^[ ]*//g')
  if [ ! -f "$t/README.md" ]; then echo "$msg" | tee "$t/README.md";
  else echo "$msg"
  fi
done
