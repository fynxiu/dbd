#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
unset CDPATH

DBD_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."

source "${DBD_ROOT}/hack/util/log.sh"
