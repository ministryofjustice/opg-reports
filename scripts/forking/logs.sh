#!/usr/bin/env bash
set -eo pipefail

log() {
    local level=${1:-INFO}
    local msg="${2}"
    local output=$(test ${LOG_LEVEL} -ge ${level} && echo "true" || echo "false")
    if [[ "${output}" == "true" ]]; then
        echo "${msg}"
    fi
    if [[ "${level}" == "${ERROR}" ]]; then
        exit 1
    fi
}
