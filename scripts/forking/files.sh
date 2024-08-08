#!/usr/bin/env bash
set -eo pipefail

delete_files() {
    local directory="${1}"
    local pattern="${2}"
    local -n exclusions="${3}"

    printf '1: %s\n' "${directory}"
    printf '2: %s\n' "${pattern}"
    printf '3: %q\n' "${exclusions[@]}"


    for file in ${directory}/${pattern}; do
        local base=$(basename "${file}")
        local should_delete="true"
        local deleted="${N}"
        local is_excluded=$(echo ${exclusions[@]} | grep -ow "${base}" | wc -w | tr -d ' ')

        # this file matches an excluded file, so skip it
        if [[ "${is_excluded}" == "1" ]]; then
            info "${base}" "delete" "${SKIP}"
            continue
        fi

        # delete the file
        LIVE && \
            rm -f "${file}" && \
            info "${base}" "delete" "${Y}"

    done

    flush
    # for file in ${directory}/${pattern}; do
    #     local base=$(basename "${file}")
    #     local should_delete="true"
    #     local deleted="${N}"

    #     # check all excluded files
    #     for exclude in "${excludeArray[@]}"; do
    #         if [[ "${base}" == "${exclude}" ]]; then
    #           should_delete="false"
    #         fi
    #     done


    #     # if [[ "${base}" == "${exclude}" ]]; then
    #     #     should_delete="false"
    #     # fi

    #     # if [[ "${should_delete}" == "true" ]]; then

    #     #     # LIVE && rm -f "${file}" && deleted="${Y}" ||
    #     #     # if [[ -f "${file}" ]]; then
    #     #     # fi

    #     # else
    #     #     deleted="${SKIP}"
    #     # fi

    #     # info "${base}" "delete?" "${should_delete}" "${deleted}"
    # done
    log ${INFO} "-"
}

delete_directory() {
    local directory="${1}"

    log ${INFO} "[deleting directory]"
    log ${DEBUG} "directory: ${directory}"
    log ${INFO} ""

    LIVE && \
        rm -Rf "${directory}" && \
        log ${INFO} "${Y} deleted: ${directory}" || \
    log ${DEBUG} "dry run - skipping"

    log ${INFO} "-"
}
