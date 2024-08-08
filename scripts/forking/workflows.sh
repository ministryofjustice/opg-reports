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

gh_workflow_replace_key(){
    local field="${1}"
    local directory="${2}"
    local workflow="${3}"
    local current="${4}"
    local replacement="${5}"
    local original="${directory}/${workflow}"
    local updated="${original}.updated"
    local base=$(basename "${updated}")

    log ${DEBUG} "directory: ${directory}"
    log ${DEBUG} "workflow: ${workflow}"
    log ${DEBUG} "current: ${current}"
    log ${DEBUG} "replacement: ${replacement}"
    log ${INFO} ""

    content=$(cat "${original}")
    # sed doesnt like the : in the arn
    echo "${content//${field}: \"${current}\"/${field}: \"${replacement}\"}" > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        log ${INFO} "${Y} updated workflow: ${workflow}" || \
    log ${INFO} "${Y} generated workflow: ${base}"

}

gh_workflow_replace_bucket() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_s3_bucket"

    log ${INFO} "[updating bucket name in workflow]"
    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"

}


gh_workflow_replace_bucket_download_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_s3_download"

    log ${INFO} "[updating bucket download role in workflow]"
    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"

}


gh_workflow_replace_bucket_upload_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_bucket_upload"

    log ${INFO} "[updating bucket upload role in workflow]"
    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"

}

gh_workflow_replace_ecr_registry_id(){
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="ecr_registry_id"

    log ${INFO} "[updating ecr push role in workflow]"
    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
}

gh_workflow_replace_ecr_push_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_ecr_login_and_push"

    log ${INFO} "[updating ecr push role in workflow]"
    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"

}
