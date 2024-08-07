#!/usr/bin/env bash
set -eo pipefail

delete_files() {
    local directory="${1}"
    local pattern="${2}"
    local exclude="${3}"

    log ${INFO} "[deleting files from directory]"
    log ${INFO} "directory: ${directory}"
    log ${INFO} "pattern: ${pattern}"
    log ${INFO} "exclude: ${exclude}"
    log ${INFO} ""

    for file in ${directory}/${pattern}; do
        local base=$(basename "${file}")
        local should_delete="true"

        if [[ "${base}" == "${exclude}" ]]; then
            should_delete="false"
        fi

        log ${DEBUG} "file: ${base}"
        log ${DEBUG} "delete? ${should_delete}"

        if [[ "${should_delete}" == "true" ]]; then
            LIVE && \
                rm -f "${file}" && \
                log ${INFO} "${Y} deleted: ${base}" || \
            log ${DEBUG} "dry run - skipping"
        fi
    done
    log ${INFO} "-"
}

delete_directory() {
    local directory="${1}"

    log ${INFO} "[deleting directory]"
    log ${INFO} "directory: ${directory}"
    log ${INFO} ""

    LIVE && \
        rm -Rf "${directory}" && \
        log ${INFO} "${Y} deleted: ${directory}" || \
    log ${DEBUG} "dry run - skipping"

    log ${INFO} "-"
}
