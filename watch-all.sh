#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

trap 'kill $(jobs -p)' EXIT

# Recompile bootstrap as necessary
fswatch bootstrap/variables.less | xargs -n1 ./customize-bootstrap.sh &

# Have hugo serve and watch
hugo -w -D -d test-public server &

wait
