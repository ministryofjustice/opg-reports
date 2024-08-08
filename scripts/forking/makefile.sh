#!/usr/bin/env bash
set -eo pipefail

makefile_replace_bucket() {
    local directory="${1}"
    local file="${2}"
    local current="${3}"
    local replacement="${4}"
    local original="${1}/${file}"
    local updated="${source}.updated"
    local base=$(basename "${updated}")

    sed "s/BUCKET ?= ${current}/BUCKET ?= ${replacement}/g" ${original} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        debug "${Y}" "updated makefile bucket" "[${file}]" || \
    debug "${SKIP}" "generated makefile example" "[${base}]"

    info "${Y}" "Replaced makefile default s3 bucket"
    divider
}


makefile_replace_aws_profile() {
    local directory="${1}"
    local file="${2}"
    local current="${3}"
    local replacement="${4}"
    local original="${1}/${file}"
    local updated="${source}.updated"
    local base=$(basename "${updated}")

    sed "s/AWS_PROFILE ?= ${current}/AWS_PROFILE ?= ${replacement}/g" ${original} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        debug "${Y}" "updated makefile aws profile" "[${file}]" || \
    debug "${SKIP}" "generated makefile example" "[${base}]"

    info "${Y}" "Replaced makefile aws profile"
    divider
}
