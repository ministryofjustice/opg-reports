#!/usr/bin/env bash
set -eo pipefail

################################################
# USER INPUTS AND THEIR DEFAULTS
# - defaults are readonly and have D_ prefix
################################################
# defaults
readonly D_UNIT="OPG"
readonly D_AWS_PROFILE="shared-development-operator"
readonly D_BUCKET_NAME_DEV="report-data-development"
readonly D_BUCKET_DOWNLOAD_ROLE_DEV="arn:aws:iam::679638075911:role/opg-reports-github-actions-s3"
readonly D_BUCKET_UPLOAD_ROLE_DEV="arn:aws:iam::679638075911:role/opg-reports-github-actions-s3"
readonly D_ECR_REGISTRY_ID="311462405659"
readonly D_ECR_PUSH_ROLE_DEV="arn:aws:iam::311462405659:role/opg-reports-github-actions-ecr-push"
readonly D_GITHUB_ORG="ministryofjustice"
readonly D_GITHUB_TEAM="opg"
# when true, will ask before executing
CONFIRM="true"
# empty ones for setup with inputs / read
UNIT=""
AWS_PROFILE=""
BUCKET_NAME_DEV=""
BUCKET_DOWNLOAD_ROLE_DEV=""
BUCKET_UPLOAD_ROLE_DEV=""
ECR_REGISTRY_ID=""
ECR_PUSH_ROLE_DEV=""
GITHUB_ORG=""
GITHUB_TEAM=""
################################################
# DIRECTORY HELPER
################################################
d() { cd "${1}"; pwd; }
################################################
# DIRECTORY PATHS
################################################
readonly SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
readonly ROOT_DIR=$(d "${SCRIPT_DIR}/../")
readonly GITHUB_DIR="${ROOT_DIR}/.github"
readonly GITHUB_WORKFLOW_DIR="${GITHUB_DIR}/workflows"
readonly TERRAFORM_DIR="${ROOT_DIR}/terraform"
readonly SERVICE_DIR="${ROOT_DIR}/servers"
readonly SERVICE_FRONT_DIR="${ROOT_DIR}/servers/front"
readonly DOCKER_DIR_API="${ROOT_DIR}/docker/api"
################################################
# FILES
################################################
readonly MAKEFILE="Makefile"
readonly GITHUB_REPORT_PATTERN="report_*.yml"
readonly GITHUB_REPORTS_TO_KEEP=( "report_repository_standards.yml" )
readonly GITHUB_REPOSITORY_REPORT="report_repository_standards.yml"
readonly GITHUB_WORKFLOW_PR="workflow_pr.yml"
readonly GITHUB_WORKFLOW_LIVE="workflow_path_to_live.yml"
readonly DOCKER_COMPOSE_FILE="docker-compose.yml"
readonly FRONT_CONFIG_FILE="config.base.json"
readonly FRONT_CONFIG_LINK="config.json"
readonly DOCKER_FILE="Dockerfile"
################################################
# OUTPUT ICONS
################################################
readonly START="⇨"
readonly END="⇦"
readonly SKIP="-"
readonly Y="✔"
readonly N="𐄂"
################################################
# FIND / REPLACE
################################################
# markers within the workflows to remove
readonly CHUNK_START="#--fork-remove-start"
readonly CHUNK_END="#--fork-remove-end"
readonly TEXT_REPLACE="#--fork-replacement"
# keys to look for in the workflows
readonly YAML_BUCKET="aws_s3_bucket"
readonly YAML_S3_DOWNLOAD="aws_role_s3_download"
readonly YAML_S3_UPLOAD="aws_role_s3_upload"
readonly YAML_ECR_REGISTRY_ID="ecr_registry_id"
readonly YAML_ECR_PUSH="aws_role_ecr_login_and_push"
readonly MAKEFILE_BUCKET="BUCKET"
readonly MAKEFILE_PROFILE="AWS_VAULT_PROFILE"
readonly DOCKER_REGISTRY="image"
readonly GH_ORG_KEY="github_org"
readonly GH_TEAM_KEY="github_team"
readonly CONFIG_UNIT="organisation"
################################################
# FUNCTIONS
################################################
# standard print
p() {
    printf "%-4s %-30s %s\n" "${1}" "${2:-}" "${3:-}"
}

############## FILES
# Delete files in the directory, exclude some files
delete_files() {
    local directory="${1}"
    local pattern="${2}"
    local -n exclusions="${3}"
    p "${START}" "Deleting files..."

    for file in ${directory}/${pattern}; do
        local base=$(basename "${file}")
        local is_excluded=$(echo ${exclusions[@]} | grep -ow "${base}" | wc -w | tr -d ' ')
        # this file matches an excluded file, so skip it
        if [[ "${is_excluded}" == "1" ]]; then
            p "${SKIP}" "excluded" "[${base}]"
            continue
        fi
        # remove the file, show error if fails
        rm -f "${file}" && p "${Y}" "deleted" "[${base}]" || \
            p "${N}" "failed to delete" "[${base}]"
    done
    p "${END}" "Deleted files"
}
# Delete the directory passed
delete_directory() {
    local directory="${1}"
    local base=$(basename "${directory}")
    p "${START}" "Deleting directory..."

    rm -Rf "${directory}" && p "${Y}" "deleted" "[${base}]" || \
        p "${N}" "failed to delete" "[${base}]"

    p "${END}" "Deleted directory"
}

link(){
    local dir="${1}"
    local source="${2}"
    local target="${3}"
    p "${START}" "Updating symlink"
    cd "${dir}"
    rm -f "./${target}"
    ln -s "./${source}" "./${target}" && p "${Y}" "symlinked" "[${target}]" || \
        p "${N}" "failed to link" "[${target}]"
    cd - 2>&1 >/dev/null
    p "${END}" "Updated symlink"
}

############## FIND / REPLACE
remove_chunk() {
    local dir="${1}"
    local file="${2}"
    local start="${3}"
    local end="${4}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Removing chunk"
    sed "/${start}/,/${end}/d" ${source} > ${destination}
    mv "${destination}" "${source}" && p "${Y}" "removed chunk" "[${base}]" || \
        p "${N}" "failed removing chunk" "[${base}]"

    p "${END}" "Removed chunk"
}

remove_text() {
    local dir="${1}"
    local file="${2}"
    local original="${3}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Removing text"
    sed "s/${original}//g" ${source} > ${destination}
    mv "${destination}" "${source}" && p "${Y}" "removed text" "[${base}]" || \
        p "${N}" "failed removing text" "[${base}]"

    p "${END}" "Removed text"
}

replace_compose_attr(){
    local dir="${1}"
    local file="${2}"
    local field="${3}"
    local original="${4}"
    local replacement="${5}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Replacing compose attribute"
    content=$(cat "${source}")
    echo "${content//${field}: ${original}/${field}: ${replacement}}" > ${destination}

    mv "${destination}" "${source}" && p "${Y}" "replaced attr" "[${base}]" || \
        p "${N}" "failed replaced attr" "[${base}]"

    p "${END}" "Replaed compose attribute"
}

replace_yaml_attr() {
    local dir="${1}"
    local file="${2}"
    local field="${3}"
    local original="${4}"
    local replacement="${5}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Replacing yaml attribute"
    content=$(cat "${source}")
    echo "${content//${field}: \"${original}\"/${field}: \"${replacement}\"}" > ${destination}

    mv "${destination}" "${source}" && p "${Y}" "replaced attr" "[${base}]" || \
        p "${N}" "failed replaced attr" "[${base}]"

    p "${END}" "Replaced yaml attribute"
}

replace_makefile_var() {
    local dir="${1}"
    local file="${2}"
    local field="${3}"
    local original="${4}"
    local replacement="${5}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Replacing makefile var"
    sed "s/${field} ?= ${original}/${field} ?= ${replacement}/g" ${source} > ${destination}

    mv "${destination}" "${source}" && p "${Y}" "replaced makefile var" "[${base}]" || \
        p "${N}" "failed replace var" "[${base}]"
    p "${END}" "Replaced makefile var"
}

replace_config_attr() {
    local dir="${1}"
    local file="${2}"
    local field="${3}"
    local original="${4}"
    local replacement="${5}"
    local source="${1}/${2}"
    local destination="${source}.copy"
    local base=$(basename "${source}")

    p "${START}" "Replacing config attr"
    content=$(cat "${source}")
    echo "${content//\"${field}\": \"${original}\"/\"${field}\": \"${replacement}\"}" > ${destination}

    mv "${destination}" "${source}" && p "${Y}" "replaced attr" "[${base}]" || \
        p "${N}" "failed replaced attr" "[${base}]"
    p "${END}" "Replaced config attr"

}
################################################
# MESSAGES
################################################
secrets() {
    echo "------------------"
    echo "!! You will need to set the following secrets on your fork !!"
    echo "GH_ORG_ACCESS_TOKEN"
    echo "This is a token with access to all public and private repositories in the org and team you want report on."
    echo "------------------"
}
################################################
# ARGUMENT HANDLING
################################################
# read in cli flags to the vars to skip the read -p calls
args() {
    while [[ "$#" -gt 0 ]]; do
        case $1 in
            # with -y set it no longer asks for confirmation
            -y) CONFIRM="false";;
            --business-unit) UNIT="${2}"; shift ;;
            --aws-profile) AWS_PROFILE="${2}"; shift ;;
            --development-bucket-name) BUCKET_NAME_DEV="${2}"; shift ;;
            --development-bucket-download-arn) BUCKET_DOWNLOAD_ROLE_DEV="${2}"; shift ;;
            --development-bucket-upload-arn) BUCKET_UPLOAD_ROLE_DEV="${2}"; shift ;;
            --ecr-registry-id) ECR_REGISTRY_ID="${2}"; shift ;;
            --ecr-login-push-arn) ECR_PUSH_ROLE_DEV="${2}"; shift ;;
            --gh-org) GITHUB_ORG="${2}"; shift ;;
            --gh-team) GITHUB_TEAM="${2}"; shift ;;
            *) echo "Unknown parameter passed: $1"; exit 1 ;;
        esac
        shift
    done
}
# ask for values when empty
reads(){
    # business unit
    if [[ "${UNIT}" == "" ]]; then
        read -p "Name of business unit: [${D_UNIT}] " UNIT
        UNIT="${UNIT:-$D_UNIT}"
    fi
    #bucket
    if [[ "${BUCKET_NAME_DEV}" == "" ]]; then
        read -p "Name of the *DEVELOPMENT* S3 bucket used for data storage: [${D_BUCKET_NAME_DEV}] " BUCKET_NAME_DEV
        BUCKET_NAME_DEV="${BUCKET_NAME_DEV:-$D_BUCKET_NAME_DEV}"
    fi
    # aws profile used for downloading to local
    if [[ "${AWS_PROFILE}" == "" ]]; then
        read -p "The *DEVELOPMENT* aws profile for local S3 download: [${D_AWS_PROFILE}] " AWS_PROFILE
        AWS_PROFILE="${AWS_PROFILE:-$D_AWS_PROFILE}"
    fi
    # download
    if [[ "${BUCKET_DOWNLOAD_ROLE_DEV}" == "" ]]; then
        read -p "The *DEVELOPMENT* OIDC role ARN to use for *DOWNLOADING* from the S3 bucket in workflows: [${D_BUCKET_DOWNLOAD_ROLE_DEV}] " BUCKET_DOWNLOAD_ROLE_DEV
        BUCKET_DOWNLOAD_ROLE_DEV="${BUCKET_DOWNLOAD_ROLE_DEV:-$D_BUCKET_DOWNLOAD_ROLE_DEV}"
    fi
    # upload
    if [[ "${BUCKET_UPLOAD_ROLE_DEV}" == "" ]]; then
        read -p "The *DEVELOPMENT* OIDC role ARN to use for *UPLOADING* to the S3 bucket in workflows: [${D_BUCKET_UPLOAD_ROLE_DEV}] " BUCKET_UPLOAD_ROLE_DEV
        BUCKET_UPLOAD_ROLE_DEV="${BUCKET_UPLOAD_ROLE_DEV:-$D_BUCKET_UPLOAD_ROLE_DEV}"
    fi
    # ecr
    if [[ "${ECR_REGISTRY_ID}" == "" ]]; then
        read -p "The AWS ECR registry id: [${D_ECR_REGISTRY_ID}] " ECR_REGISTRY_ID
        ECR_REGISTRY_ID="${ECR_REGISTRY_ID:-$D_ECR_REGISTRY_ID}"
    fi
    # ecr role
    if [[ "${ECR_PUSH_ROLE_DEV}" == "" ]]; then
        read -p "The *DEVELOPMENT* OIDC role ARN to use for pushing to ECR: [${D_ECR_PUSH_ROLE_DEV}] " ECR_PUSH_ROLE_DEV
        ECR_PUSH_ROLE_DEV="${ECR_PUSH_ROLE_DEV:-$D_ECR_PUSH_ROLE_DEV}"
    fi

    if [[ "${GITHUB_ORG}" == "" ]]; then
        read -p "The github organisation slug: [${D_GITHUB_ORG}] " GITHUB_ORG
        GITHUB_ORG="${GITHUB_ORG:-$D_GITHUB_ORG}"
    fi
    if [[ "${GITHUB_TEAM}" == "" ]]; then
        read -p "The github parent teams slug: [${D_GITHUB_TEAM}] " GITHUB_TEAM
        GITHUB_TEAM="${GITHUB_TEAM:-$D_GITHUB_TEAM}"
    fi
}

# show the set values
show_conf(){
    echo "
------------------
CONFIG:"

echo "
Business unit name:@[${UNIT}]@
AWS S3 nucket name:@[${BUCKET_NAME_DEV}]
AWS profile for *LOCAL* s3 bucket *DOWNLOAD*:@[${AWS_PROFILE}]
AWS role arn for *WORKFLOW* s3 bucket *DOWNLOAD*:@[${BUCKET_DOWNLOAD_ROLE_DEV}]
AWS role arn for *WORKFLOW* s3 bucket *UPLOAD*:@[${BUCKET_UPLOAD_ROLE_DEV}]
AWS ECR registry id:@[${ECR_REGISTRY_ID}]
AWS OIDC role arn for the *WORKFLOW* ecr login and push:@[${ECR_PUSH_ROLE_DEV}]
GITHUB organisation:@[${GITHUB_ORG}]
GITHUB parent team:@[${GITHUB_TEAM}]" | column -s@ -t
echo ""
}
################################################
# PROCESS CONTROL
################################################
# ask to continue or not
should_continue() {
    local continue="Y"
    if [[ "${CONFIRM}" == "true" ]]; then
        read -p "Contine using the above details? (Y|N) [${continue}] " continue
        echo ""
        continue="${continue:-"Y"}"
        if [[ "${continue}" != [Yy] ]]; then
            echo "exiting..."
            exit 1
        fi
    fi

}
################################################
# MAIN
################################################
main(){
    # setup calls
    args "$@"
    reads
    show_conf
    should_continue

    # delete very specifc files
    rm -f "${ROOT_DIR}/scripts/create-local-aws-costs-script"
    # delete unused files
    delete_files "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPORT_PATTERN}" GITHUB_REPORTS_TO_KEEP
    # remove terraform directory
    delete_directory "${TERRAFORM_DIR}"
    # swap config files to simple version & replace org
    link "${SERVICE_FRONT_DIR}" "${FRONT_CONFIG_FILE}" "${FRONT_CONFIG_LINK}"
    replace_config_attr "${SERVICE_FRONT_DIR}" "${FRONT_CONFIG_FILE}" "${CONFIG_UNIT}" "${D_UNIT}" "${UNIT}"

    ############## DEVELOPMENT
    # remove terraform chunks from workflows
    remove_chunk "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${CHUNK_START}" "${CHUNK_END}"
    remove_text "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${TEXT_REPLACE}"
    # replace bucket & ecr properties
    replace_makefile_var "${ROOT_DIR}" "${MAKEFILE}" "${MAKEFILE_BUCKET}" "${D_BUCKET_NAME_DEV}" "${BUCKET_NAME_DEV}"
    replace_makefile_var "${ROOT_DIR}" "${MAKEFILE}" "${MAKEFILE_PROFILE}" "${D_AWS_PROFILE}" "${AWS_PROFILE}"
    # workflows
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${YAML_BUCKET}" "${D_BUCKET_NAME_DEV}" "${BUCKET_NAME_DEV}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${YAML_S3_DOWNLOAD}" "${D_BUCKET_DOWNLOAD_ROLE_DEV}" "${BUCKET_DOWNLOAD_ROLE_DEV}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${YAML_ECR_REGISTRY_ID}" "${D_ECR_REGISTRY_ID}" "${ECR_REGISTRY_ID}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${YAML_ECR_PUSH}" "${D_ECR_PUSH_ROLE_DEV}" "${ECR_PUSH_ROLE_DEV}"
    # replace in docker compose
    replace_compose_attr "${ROOT_DIR}" "${DOCKER_COMPOSE_FILE}" "${DOCKER_REGISTRY}" "${D_ECR_REGISTRY_ID}" "${ECR_REGISTRY_ID}"

    ############## PRODUCTION
    # remove terraform chunks from workflows
    remove_chunk "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${CHUNK_START}" "${CHUNK_END}"
    remove_text "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${TEXT_REPLACE}"

    # replace bucket & ecr properties
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${YAML_BUCKET}" "${D_BUCKET_NAME_DEV}" "${BUCKET_NAME_DEV}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${YAML_S3_DOWNLOAD}" "${D_BUCKET_DOWNLOAD_ROLE_DEV}" "${BUCKET_DOWNLOAD_ROLE_DEV}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${YAML_ECR_REGISTRY_ID}" "${D_ECR_REGISTRY_ID}" "${ECR_REGISTRY_ID}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_LIVE}" "${YAML_ECR_PUSH}" "${D_ECR_PUSH_ROLE_DEV}" "${ECR_PUSH_ROLE_DEV}"

    # main report - swap props
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPOSITORY_REPORT}" "${YAML_BUCKET}" "${D_BUCKET_NAME_DEV}" "${BUCKET_NAME_DEV}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPOSITORY_REPORT}" "${YAML_S3_UPLOAD}" "${D_BUCKET_UPLOAD_ROLE_DEV}" "${BUCKET_UPLOAD_ROLE_DEV}"

    # swap the org & team for the workflow run
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPOSITORY_REPORT}" "${GH_ORG_KEY}" "${D_GITHUB_ORG}" "${GITHUB_ORG}"
    replace_yaml_attr "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPOSITORY_REPORT}" "${GH_TEAM_KEY}" "${D_GITHUB_TEAM}" "${GITHUB_TEAM}"

    # remove chunks from makefile
    remove_chunk "${ROOT_DIR}" "${MAKEFILE}" "${CHUNK_START}" "${CHUNK_END}"
    # remove chunks from dockerfile
    remove_chunk "${DOCKER_DIR_API}" "${DOCKER_FILE}" "${CHUNK_START}" "${CHUNK_END}"

    secrets
}



main "$@"
