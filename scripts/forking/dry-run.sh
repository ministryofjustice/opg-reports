#!/usr/bin/env bash
set -eo pipefail

# when DRY_RUN is enabled (set as 1), then LIVE
# returns false (1)
# when dry run is disabled (set as 0), then LIVE
# returns true (0)
# defaults to false
LIVE() {
    if [[ "${DRY_RUN}" == "${TRUE_FLAG}" ]]; then
        return 1
    elif [[ "${DRY_RUN}" == "${FALSE_FLAG}" ]]; then
        return 0
    fi
    return 1
}
