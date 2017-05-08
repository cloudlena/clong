#!/bin/bash

set -e -u

if [ -z "${PORT}" ]; then
  echo "Error: No PORT found" >&2
  exit 1
fi

# Run binary
./clong -port "${PORT}"
