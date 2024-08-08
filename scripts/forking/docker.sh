#!/usr/bin/env bash
set -eo pipefail

docker_compose_replace_registry() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="image"
    local original="${directory}/${workflow}"
    local updated="${original}.updated"
    local base=$(basename "${updated}")

    content=$(cat "${original}")
    echo "${content//${field}: ${current}/${field}: ${replacement}}" > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        log ${INFO} "${Y} updated workflow: ${workflow}" || \
    log ${INFO} "${Y} generated workflow: ${base}"

}
