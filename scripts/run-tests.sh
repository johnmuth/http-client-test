#!/bin/bash -e

set -o errexit
set -o nounset
set -o pipefail
set -e
set -u

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

source "$SCRIPT_DIR/build-functions.sh"

fmtAppCode
lintAppCode
runTests
