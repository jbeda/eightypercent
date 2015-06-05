#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

trap 'kill $(jobs -p)' EXIT

# Recompile bootstrap as necessary
fswatch-run bootstrap/variables.less ./customize-bootstrap.sh &

# Have hugo serve and watch
hugo -w -d test-public server &

wait