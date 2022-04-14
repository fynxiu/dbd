#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

DBD_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
source "${DBD_ROOT}/hack/vcs.sh"
source "${DBD_ROOT}/hack/version.sh"


readonly CURRENT_VERSION=$(dbd::vcs::version)
readonly NEW_VERSION=$(dbd::version::next_version "${CURRENT_VERSION}" ${1:-patch} ${2:-} ${3:-})
if [[ ${NEW_VERSION} == "" ]] ; then
  exit 1
fi
if ! dbd::version::validate_version "${NEW_VERSION}"; then
  echo "invalid new version : ${NEW_VERSION}" >&2
  exit 1
fi

readonly VERSION_FILE="${DBD_ROOT}/internal/version/version.g.go"
echo -e "// GENERATED CODE by dbd. DO NOT EDIT.\n\n" \
  "package version\n\n" \
  "// Version is the version of dbd.\n" \
  "const Version = \"${NEW_VERSION}\"\n" > ${VERSION_FILE}
gofmt -w -s "${VERSION_FILE}"

# git commit and tag
git add "${VERSION_FILE}"
if [[ $(git tag -l "${NEW_VERSION}") ]]; then
  log::fatal "tag ${NEW_VERSION} already exists"
fi
git commit -m "Update version to ${NEW_VERSION}"
git tag -a "${NEW_VERSION}" -m "Version ${NEW_VERSION}"
