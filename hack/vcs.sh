#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

# dbd::vcs::version - Print the version of the current VCS
function dbd::vcs::version {
    local version
    version=$(command git describe --tags --abbrev=0 HEAD 2>/dev/null) \
    || version="v1.0.0-alpha.1"
    echo "${version}"
}

# dbd::vcs::commit_hash - Print the current VCS commit short hash
function dbd::vcs::commit_hash {
  echo $(command git rev-parse --short HEAD 2>/dev/null)
}
