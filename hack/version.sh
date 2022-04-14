#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

DBD_ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
source "${DBD_ROOT}/hack/init.sh"

readonly SEMVER_REGEX="^([vV]?)(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z]+(\\.[0-9A-Za-z]+)*)?(\\+[0-9A-Za-z]+(\\.[0-9A-Za-z]+)*)?$"

function dbd::version::next_version {
  local current_version=$1
  local option=$2
  local prere=${3:-}
  local build=${4:-}
  [[ ${prere} == "" ]] || [[ ${prere} == -* ]] || prere="-"${prere}
  [[ ${build} == "" ]] || [[ ${build} == +* ]] || build="+"${build}

  if [[ ! (${current_version} =~ ${SEMVER_REGEX}) ]]; then
    echo "invalid current version : ${current_version}" >&2
    return 0
  else
    local prefix=${BASH_REMATCH[1]}
    local major=${BASH_REMATCH[2]}
    local minor=${BASH_REMATCH[3]}
    local patch=${BASH_REMATCH[4]}
    if [[ ${option} == "build" ]]; then
      if [[ ${build} == "" || ${prere} == "" ]]; then
        echo "prere and build must be provided" >&2
        return 0
      fi
    fi

    readonly new_version=${prefix}$(
      case "${option}" in
      major) echo "$((major + 1)).0.0" ;;
      minor) echo "$major.$((minor + 1)).0" ;;
      prere) echo "$major.$minor.$patch" ;;
      build) echo "$major.$minor.$patch" ;;
      patch | *) echo "$major.$minor.$((patch + 1))" ;;
      esac
    )$prere$build

    echo ${new_version}
  fi

}

function dbd::version::validate_version {
  [[ "${1}" =~ ${SEMVER_REGEX} ]]
}

# readonly xxx=$(dbd::version::next_version "v1.2.3-alpha.1" patch alpha.2 build.4)
# if [[ $xxx != "" ]]; then
#   echo "${xxx}"
#   dbd::version::validate_version "$xxx" || log::fatal "invalid new version : ${xxx}"
#   echo "new version: ${xxx}"
# else
#   echo "failed"
# fi
