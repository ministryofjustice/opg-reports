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
readonly GITHUB_REPORT_KEEP="report_repository_standards"
################################################

# set the default log level to info
LOG_LEVEL=${DEBUG}
# set if dry run or not - default to true
DRY_RUN=${FALSE_FLAG}
