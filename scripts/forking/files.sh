#!/usr/bin/env bash
set -eo pipefail

delete_files() {
    local directory="${1}"
    local pattern="${2}"
    local -n exclusions="${3}"
    local live=$(LIVE && echo "true" || echo "false")

    for file in ${directory}/${pattern}; do
        local base=$(basename "${file}")
        local should_delete="true"
        local deleted="${N}"
        local is_excluded=$(echo ${exclusions[@]} | grep -ow "${base}" | wc -w | tr -d ' ')

        # this file matches an excluded file, so skip it
        if [[ "${is_excluded}" == "1" ]]; then
            debug "${SKIP}" "delete" "[${base}]"
            continue
        fi
        # if its a dry run, skip the file delete
        if [[ "${live}" != "true" ]]; then
            debug "${SKIP}" "delete" "[${base}]"
            continue
        fi
        # delete the file
        LIVE && \
            rm -f "${file}" && \
            debug "${Y}" "delete" "[${base}]" || \
        err "${N}" "delete" "[${base}]"

    done

    info "${Y}" "Deleted files"
    divider
}

delete_directory() {
    local directory="${1}"
    local live=$(LIVE && echo "true" || echo "false")

    # if this isnt a live run, then output info
    if [[ "${live}" != "true" ]]; then
        debug "${SKIP}" "delete" "[${directory}]"
    # if its live and works, output info, otherwise flag error
    else
        rm -Rf "${directory}" && \
            debug "${Y}" "delete" "[${directory}]" || \
            err "${N}" "delete" "[${directory}]"
    fi

    info "${Y}" "Deleted directory"
    divider
}
