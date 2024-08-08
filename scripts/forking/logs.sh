#!/usr/bin/env bash
set -eo pipefail

################################################
_buffer_info=""
################################################

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


debug() {
    echo ""
}

info() {

    for i in "$@"; do
        local l="${#i}"
        printf -v ${_buffer_info} "%-${l}s   " "${i}"
    done
    printf -v ${_buffer_info} "\n"

}

err() {

}


flush() {
    echo "${_buffer_info}" | column -s: -t
}
