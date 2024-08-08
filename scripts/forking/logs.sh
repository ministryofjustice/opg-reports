#!/usr/bin/env bash
set -eo pipefail

################################################
buffer=""
################################################

divider() {
    local dev="------------------------------"
    info "${dev}"
    flush
}

err(){
    info "$@"
}

debug() {
    info "$@"
}

info() {
    local n=$'\n'

    for i in "$@"; do
        local buff=""
        printf -v buff "%s\t" "${i}"
        buffer="${buffer}${buff}"
    done
    buffer="${buffer}${n}"

}



flush() {
    echo "${buffer}" | column -s: -t
    buffer=""
}
