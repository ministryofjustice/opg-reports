#!/usr/bin/env bash
set -eo pipefail

gh_workflow_remove_marked() {
    local directory="${1}"
    local workflow="${2}"
    local original="${directory}/${workflow}"
    local remove="${original}.removed"
    local updated="${original}.updated"
    local base=$(basename "${cpy}")

    sed '/#--fork-remove-start/,/#--fork-remove-end/d' ${original} > ${remove}
    sed 's/#--fork-replacement//g' ${remove} > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${remove}" && \
        debug "${Y}" "removed marked segments" "[${workflow}]" || \
    debug "${SKIP}" "generated example workflows" "[${base}]"

    info "${Y}" "Removed workflow segments"
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

    content=$(cat "${original}")
    # sed doesnt like the : in the arn
    echo "${content//${field}: \"${current}\"/${field}: \"${replacement}\"}" > ${updated}

    LIVE && \
        mv "${updated}" "${original}" && \
        rm -f "${updated}" && \
        debug "${Y}" "updated workflow" "[${workflow}]" || \
    debug "${SKIP}" "generated workflow" "[${base}]"

}

gh_workflow_replace_bucket() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_s3_bucket"

    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
    info "${Y}" "Replaced s3 bucket name"
}


gh_workflow_replace_bucket_download_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_s3_download"

    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
    info "${Y}" "Replaced s3 download role"

}


gh_workflow_replace_bucket_upload_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_bucket_upload"

    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
    info "${Y}" "Replaced s3 upload role"

}

gh_workflow_replace_ecr_registry_id(){
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="ecr_registry_id"

    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
    info "${Y}" "Replaced ecr registry id"
}

gh_workflow_replace_ecr_push_role() {
    local directory="${1}"
    local workflow="${2}"
    local current="${3}"
    local replacement="${4}"
    local field="aws_role_ecr_login_and_push"

    gh_workflow_replace_key "${field}" "${directory}" "${workflow}" "${current}" "${replacement}"
    info "${Y}" "Replaced ecr login and push role"
}
