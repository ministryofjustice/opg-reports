#!/usr/bin/env bash
set -eo pipefail

################################################
buffer=""
################################################

divider() {
    local dev="------------------------------"
    echo "${dev}"
}

err(){
    info "$@"
    flush
    exit 1
}

debug() {
    info "$@"
}

info() {
    local n=$'\n'

    for i in "$@"; do
        local buff=""
        printf -v buff "%s" "${i}"
        buffer=$(echo "${buffer}${buff}@" | tr -s '[:space:]')
    done
    buffer="${buffer}${n}"

}



flush() {
    echo "${buffer}"| column -c 3 -s@ -t
    buffer=""
}
