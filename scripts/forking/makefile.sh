#!/usr/bin/env bash
set -eo pipefail

makefile_replace_bucket() {
    local directory="${1}"
    local file="${2}"
    local original="${1}/${file}"
    local updated="${source}.updated"
    local current="${3}"
    local replacement="${4}"

    log ${INFO} "[updating bucket name in makefile]"
    log ${DEBUG} "directory: ${directory}"
    log ${DEBUG} "makefile: ${workflow}"
    log ${DEBUG} "current: ${current}"
    log ${DEBUG} "replacement: ${replacement}"
    log ${INFO} ""

    sed "s/BUCKET ?= ${current}/BUCKET ?= ${replacement}/g" ${original} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        log ${INFO} "${Y} updated bucket: ${workflow}" || \
    log ${INFO} "${Y} generated example makefile"
}
