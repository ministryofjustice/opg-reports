#!/usr/bin/env bash
set -eo pipefail

gh_workflow_remove_marked() {
    local directory="${1}"
    local workflow="${2}"
    local source="${1}/${2}"
    local remove="${source}.remove"
    local replace="${source}.updated"
    local base=$(basename "${cpy}")

    log ${INFO} "[removing marked items from workflow]"
    log ${INFO} "directory: ${directory}"
    log ${INFO} "workflow: ${workflow}"
    log ${INFO} ""

    sed '/#--fork-remove-start/,/#--fork-remove-end/d' ${source} > ${remove}
    sed 's/#--fork-replacement//g' ${remove} > ${replace}

    LIVE && \
        mv "${replace}" "${source}" && \
        rm -f "${remove}" && \
        log ${INFO} "${Y} removed marked: ${workflow}" || \
    log ${INFO} "${Y} generated example workflows"

    log ${INFO} "-"
}
