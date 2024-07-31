#!/usr/bin/env bash


SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BASE_DIR="${SCRIPT_DIR}/../"
DOCS_DIR="${BASE_DIR}docs/go/"


rm -Rf ${DOCS_DIR}
mkdir -p ${DOCS_DIR}

for d in $(find . -type f -name '*.go' | sed -r 's|/[^/]+$||' |sort -u); do
    action="processing"
    status="✅"
    # exclude the root and terraform folders
    if [[ "${d}" != "." && "${d}" != *"terraform"* ]]; then
        dir="${DOCS_DIR}${d}"
        mkdir -p ${dir}
        go doc -all -u ${d} > ${dir}/index.txt || status="❌"
    else
        action="skipping"
        status="⛔"
    fi
    printf '%-64s %-15s %-6s\n' "${d}" "${action}" "${status}"
done
