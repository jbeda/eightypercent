#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

ROOT=$(cd $(dirname "${BASH_SOURCE}"); pwd)
BOOTSTRAP="${ROOT}/../bootstrap"

cp "${ROOT}/bootstrap/variables.less" "${BOOTSTRAP}/less/variables.less"

cd "${BOOTSTRAP}"

grunt dist

cp -vR "${BOOTSTRAP}/dist/" "${ROOT}/static/"
