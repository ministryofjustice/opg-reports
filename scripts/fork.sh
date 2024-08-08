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
readonly D_BUCKET_DOWNLOAD_ROLE_DEV="arn:aws:iam::679638075911:role/docs-and-metadata-ci"
readonly D_BUCKET_UPLOAD_ROLE_DEV="arn:aws:iam::679638075911:role/opg-reports-github-actions-s3"
readonly D_ECR_REGISTRY_ID="311462405659"
readonly D_ECR_PUSH_ROLE_DEV="arn:aws:iam::311462405659:role/opg-reports-github-actions-ecr-push"
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
readonly CHUNK_START="#--fork-remove-start"
readonly CHUNK_END="#--fork-remove-end"
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

        rm -f "${file}" && p "${Y}" "deleted" "[${base}]" || \
            p "${N}" "failed to delete" "[${base}]"
    done
    p "${END}" "Deleted files"
}
# Delete the directory passed
delete_directory() {
    local directory="${1}"

    p "${START}" "Deleting directory..."

    rm -Rf "${directory}" && p "${Y}" "deleted" "[${directory}]" || \
        p "${N}" "failed to delete" "[${directory}]"

    p "${END}" "Deleted directory"
}

############## DELETE CHUNKS
remove_chunk() {
    local dir="${1}"
    local file="${2}"
    local start="${3}"
    local end="${4}"
    local original="${1}/${2}"
    local copy="${original}.copy"

    sed "/${start}/,/${end}/d" ${original} > ${copy}
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
        read -p "The *DEVELOPMENT* role ARN to use for *DOWNLOADING* from the S3 bucket in workflow: [${D_BUCKET_DOWNLOAD_ROLE_DEV}] " BUCKET_DOWNLOAD_ROLE_DEV
        BUCKET_DOWNLOAD_ROLE_DEV="${BUCKET_DOWNLOAD_ROLE_DEV:-$D_BUCKET_DOWNLOAD_ROLE_DEV}"
    fi
    # upload
    if [[ "${BUCKET_UPLOAD_ROLE_DEV}" == "" ]]; then
        read -p "The *DEVELOPMENT* OIDC role ARN to use for *UPLOADING* to the S3 bucket in workflow: [${D_BUCKET_UPLOAD_ROLE_DEV}] " BUCKET_UPLOAD_ROLE_DEV
        BUCKET_UPLOAD_ROLE_DEV="${BUCKET_UPLOAD_ROLE_DEV:-$D_BUCKET_UPLOAD_ROLE_DEV}"
    fi
    # ecr
    if [[ "${ECR_REGISTRY_ID}" == "" ]]; then
        read -p "The AWS ECR registry id: [${D_ECR_REGISTRY_ID}] " ECR_REGISTRY_ID
        ECR_REGISTRY_ID="${ECR_REGISTRY_ID:-$D_ECR_REGISTRY_ID}"
    fi
    # ecr
    if [[ "${ECR_PUSH_ROLE_DEV}" == "" ]]; then
        read -p "The *DEVELOPMENT* OIDC role ARN to use for pushing to ECR: [${D_ECR_PUSH_ROLE_DEV}] " ECR_PUSH_ROLE_DEV
        ECR_PUSH_ROLE_DEV="${ECR_PUSH_ROLE_DEV:-$D_ECR_PUSH_ROLE_DEV}"
    fi
}

# show the set values
show_conf(){
    echo "CONFIG:
Business unit name:@[${UNIT}]@
AWS S3 nucket name:@[${BUCKET_NAME_DEV}]
AWS profile for *LOCAL* s3 bucket *DOWNLOAD*:@[${AWS_PROFILE}]
AWS role arn for *WORKFLOW* s3 bucket *DOWNLOAD*:@[${BUCKET_DOWNLOAD_ROLE_DEV}]
AWS role arn for *WORKFLOW* s3 bucket *UPLOAD*:@[${BUCKET_UPLOAD_ROLE_DEV}]
AWS ECR registry id:@[${ECR_REGISTRY_ID}]
AWS OIDC role arn for the *WORKFLOW* ecr login and push:@[${ECR_PUSH_ROLE_DEV}]" | column -s@ -t
}
################################################
# PROCESS CONTROL
################################################
# ask to continue or not
should_continue() {
    local continue="Y"
    if [[ "${CONFIRM}" == "true" ]]; then
        read -p "Contine using the above details? (Y|N) [${continue}] " continue
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
    # args "$@"
    # reads
    # show_conf
    # should_continue

    # delete unused files
    delete_files "${GITHUB_WORKFLOW_DIR}" "${GITHUB_REPORT_PATTERN}" GITHUB_REPORTS_TO_KEEP
    # remove terraform directory
    # delete_directory "${TERRAFORM_DIR}"

    # remove terraform chunks from workflows
    remove_chunk "${GITHUB_WORKFLOW_DIR}" "${GITHUB_WORKFLOW_PR}" "${CHUNK_START}" "${CHUNK_END}"
}



main "$@"
