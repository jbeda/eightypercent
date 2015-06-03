#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# Recompile bootstrap as necessary
fswatch-run bootstrap/variables.less ./customize-bootstrap.sh &

# Have hugo serve and watch
hugo -w server &
