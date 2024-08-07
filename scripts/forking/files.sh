#!/usr/bin/env bash
set -eo pipefail

delete_from_directory() {
    local directory="${1}"
    local pattern="${2}"
    local exclude="${3}"

    log ${INFO} "[deleting files from directory]"
    log ${INFO} "directory: ${directory}"
    log ${INFO} "pattern: ${pattern}"
    log ${INFO} "exclude: ${exclude}"

    for file in ${directory}/${pattern}; do
        local base=$(basename "${file}")
        local should_delete="true"

        if [[ "${base}" == "${exclude}" ]]; then
            should_delete="false"
        fi

        log ${INFO} "file: ${base}"
        log ${INFO} "delete? ${should_delete}"

        if [[ "${should_delete}" == "true" ]]; then
            LIVE() && rm -f "${file}" && log ${INFO} "${Y} deleted" || log ${DEBUG} "dry run - skipping"
        fi


    done

}
