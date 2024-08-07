#!/usr/bin/env bash
set -eo pipefail

gh_workflow_remove_marked() {
    local directory="${1}"
    local workflow="${2}"
    local original="${directory}/${workflow}"
    local remove="${original}.removed"
    local updated="${original}.updated"
    local base=$(basename "${cpy}")

    log ${INFO} "[removing marked items from workflow]"
    log ${DEBUG} "directory: ${directory}"
    log ${DEBUG} "workflow: ${workflow}"
    log ${INFO} ""

    sed '/#--fork-remove-start/,/#--fork-remove-end/d' ${original} > ${remove}
    sed 's/#--fork-replacement//g' ${remove} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${remove}" && \
        log ${INFO} "${Y} removed marked: ${workflow}" || \
    log ${INFO} "${Y} generated example workflows"

    log ${INFO} "-"
}


gh_workflow_replace_bucket() {
    local directory="${1}"
    local workflow="${2}"
    local original="${directory}/${workflow}"
    local updated="${original}.updated"
    local current="${3}"
    local replacement="${4}"

    log ${INFO} "[updating bucket name in workflow]"
    log ${DEBUG} "directory: ${directory}"
    log ${DEBUG} "workflow: ${workflow}"
    log ${DEBUG} "current: ${current}"
    log ${DEBUG} "replacement: ${replacement}"
    log ${INFO} ""

    sed "s/aws_s3_bucket: \"${current}\"/aws_s3_bucket: \"${replacement}\"/g" ${original} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        log ${INFO} "${Y} updated bucket: ${workflow}" || \
    log ${INFO} "${Y} generated example workflow"

}
