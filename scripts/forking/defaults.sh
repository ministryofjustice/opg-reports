#!/usr/bin/env bash
set -eo pipefail

################################################
readonly ERROR=1
readonly INFO=2
readonly DEBUG=3

readonly TRUE_FLAG=1
readonly FALSE_FLAG=0

readonly Y="✅"
readonly N="❌"

readonly GITHUB_DIR="${ROOT_DIR}/.github"
readonly GITHUB_WORKFLOW_DIR="${GITHUB_DIR}/workflows"
readonly GITHUB_REPORT_PATTERN="report_*.yml"
readonly GITHUB_REPORT_KEEP="report_repository_standards.yml"
readonly GITHUB_WORKFLOW_PR="workflow_pr.yml"
readonly GITHUB_WORKFLOW_LIVE="workflow_path_to_live.yml"

readonly TERRAFORM_DIR="${ROOT_DIR}/terraform"

readonly MAKEFILE="Makefile"
readonly MAKEFILE_AWS_PROFILE="shared-development-operator"

readonly BUCKET_NAME_DEV="report-data-development"
readonly BUCKET_DOWNLOAD_ROLE_DEV="arn:aws:iam::679638075911:role/docs-and-metadata-ci"
readonly BUCKET_DOWNLOAD_ROLE_PROD="arn:aws:iam::679638075911:role/docs-and-metadata-ci"
################################################

# set the default log level to info
LOG_LEVEL=${DEBUG}
# set if dry run or not - default to true
DRY_RUN=${FALSE_FLAG}
# DRY_RUN=${TRUE_FLAG}
