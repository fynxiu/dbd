#!/usr/bin/env bash

# Based on https://github.com/golang/pkgsite/commit/29728f9d1dc17d2bc77e401b00bca0cb563117fc
RED=; GREEN=; YELLOW=; NORMAL=;
MAXWIDTH=0

if tput setaf 1 >& /dev/null; then
  RED=`tput setaf 1`
  GREEN=`tput setaf 2`
  YELLOW=`tput setaf 3`
  NORMAL=`tput sgr0`
  MAXWIDTH=$(( $(tput cols) - 2 ))
fi

EXIT_CODE=0
function log::info() { echo "${GREEN}$@${NORMAL}" 1>&2; }
function log::warn() { echo "${YELLOW}$@${NORMAL}" 1>&2; }
function log::error() { echo -e "${RED}$@${NORMAL}" 1>&2; EXIT_CODE=1; }
function log::fatal() { log::error $@; exit 1; }
