#!/usr/bin/env bash
set -eo pipefail

docker_compose_replace_registry() {
    local directory="${1}"
    local file="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="image"
    local original="${directory}/${file}"
    local updated="${original}.updated"
    local base=$(basename "${updated}")

    content=$(cat "${original}")
    echo "${content//${field}: ${current}/${field}: ${replacement}}" > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        debug "${Y}" "updated docker compose" "${file}" || \
    debug "${SKIP}" "generated example docker compose" "${base}"

    info "${Y}" "Replaced docker compose registry id"
    divider

}
