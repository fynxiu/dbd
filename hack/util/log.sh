#!/usr/bin/env bash

function log::fatal() {
  echo "fatal: $*" 1>&2
  exit 1
}

function log::warn() {
  echo "warning: $*" 1>&2
}

function log::info() {
  echo "info: $*" 1>&2
}
