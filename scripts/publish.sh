#!/bin/bash
### USAGE: $0 [-h|--help] [--dry-run] -v|--version VERSION
### Publish a new version of the module.
### Example: $0 --dry-run -v 1.2.3
### Example: $0 --dry-run -v v1.2.3-beta.1

# shellcheck disable=SC2016
usage() { cat "$0" | grep '^### ' | sed 's/^### //g; s#\$0#'"$0"'#g'; }
# see https://go.dev/doc/modules/publishing
set -euo pipefail

VERSION="${VERSION:-}"
dry_run=false
while [ -n "${1:-}" ]; do
  # shift
  case "$1" in
    -h|--help) usage; exit 0 ;;
    --dry-run) dry_run=true; shift ;;
    -v|--version) VERSION="$2"; shift; shift;;
    -v=*|--version=*) VERSION="${1%=*}"; shift ;;
    *) echo "Unknown argument: $1" >&2; usage; exit 1 ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  usage
  exit 1
fi

case "$VERSION" in
  [0-9].[0-9].[0-9]         | \
  [0-9].[0-9].[0-9]-alpha   | \
  [0-9].[0-9].[0-9]-alpha.* | \
  [0-9].[0-9].[0-9]-beta    | \
  [0-9].[0-9].[0-9]-beta.*  )
    VERSION="v$VERSION"
    ;;
  v[0-9].[0-9].[0-9]         | \
  v[0-9].[0-9].[0-9]-alpha   | \
  v[0-9].[0-9].[0-9]-alpha.* | \
  v[0-9].[0-9].[0-9]-beta    | \
  v[0-9].[0-9].[0-9]-beta.*  ) ;;
  *) echo "VERSION must be in the form vMAJOR.MINOR.PATCH, got '$VERSION'" >&2; exit 1 ;;
esac
mod=""
if [ -f go.mod ]; then
  mod=$(head -1 go.mod | cut -d' ' -f2)
fi

cmd="
  git tag $VERSION &&
  git push origin $VERSION &&
  GOPROXY=proxy.golang.org go list -m $mod@$VERSION
"

if [ "$dry_run" = true ]; then
  echo "DRY RUN: $cmd"
else
  set -x
  eval "$cmd"
  set +x
fi
