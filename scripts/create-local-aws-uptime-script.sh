#!/usr/bin/env bash
# This script generates another script that contains a series of bash commands to
# fetch historical uptime data from our accounts
#
# Fetches the accounts.aws.uptime.json data from the latest opg-metadata release, and uses
# jq and sed to parse and convert the content into command to run
#
# Variable naming
# - variables injected into the built bash script are lower case
# - variables used only for this script are upper case
# - locals in this script are lower case as well
set -eo pipefail

now=$(date +%Y-%m-%d)
end_date=$(date -j -v-1d -f "%Y-%m-%d" ${now} +%Y-%m-%d)
start_date=$(date -j -v-14d -f "%Y-%m-%d" ${end_date} +%Y-%m-%d)
bucket_path="aws/uptime/daily"
profile="shared-development-operator"
bucket="report-data-development"

#
OPEN_VSCODE=${1:-n}
SUBDIR="aws-uptime-daily"
TARGET_NAME="run"
BINARY="aws_uptime_daily"
SOURCE_FILE="accounts.aws.uptime.json"
# OS info
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/')
BUILD_ARCH="${OS}_${ARCH}"
# directory / file locations
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BASE="${SCRIPT_DIR}/../builds/scripts"
BUILD_DIR="${SCRIPT_DIR}/../builds/${BUILD_ARCH}"
WORKING_DIR="${BASE}/${SUBDIR}"
RELEASE_DIR="${WORKING_DIR}/releases"
DATA_SOURCE="${RELEASE_DIR}/${SOURCE_FILE}"
TARGET_FILE="${WORKING_DIR}/${TARGET_NAME}"
# Download latest release into local directory
function download_realease() {
    local working_dir="${WORKING_DIR}"
    local release_dir="${RELEASE_DIR}"

    rm -Rf ${release_dir}
    mkdir -p ${release_dir}
    cd ${release_dir}

    gh release download --clobber --repo ministryofjustice/opg-metadata --pattern "*.tar.gz"
    tar -xzf metadata.tar.gz

    cd ${working_dir}

    echo "✅ Downloaded and extracted latest opg-metadata release"
    echo "  - [${release_dir}]"
}

# Create the header of the generated bash file
function generate_file_head() {
    # shebang
    echo '#!/usr/bin/env bash'
    # file variables used in other segments
    echo "
# variables
start_date=\"${start_date}\"
end_date=\"${end_date}\"
profile=\"${profile}\"
bucket=\"${bucket}\"
bucket_path=\"${bucket_path}\"
binary_file=\"${BINARY}\"

rm -Rf ./data/
"
}

# generate content that asks for user input for the start & end date
function generate_file_date_input() {
    echo '
# ask for start and end dates
read -p "Start date [${start_date}]: " start
start="${start:-$start_date}"

read -p "End date [${end_date}]: " end
end="${end:-$end_date}"
'
}

# generate content that starts the looping over dates
function generate_file_date_loop_start(){
    echo '
# loop over date range provided
day="${start}"
while [ "${day}" != "${end}" ]; do
  echo "day: ${day}"
'

}

# loop end and iterate month
function generate_file_date_loop_end() {
    echo '  day=$(date -j -v +1d -f "%Y-%m-%d" ${day} +%Y-%m-%d)'
    echo 'done'

}

########
# use jq to create a string of commands from the source file
#
# due to fun with quotes and escaping we use some subs with sed
# fixing it after jq. Details in order:
#  - replace ~ with \: so the command is neatly split over multiple lines
#  - replace | with ": so argument values have string encapsulation for spaces etc
#  - replace -null- with -: some accounts dont have an environment value, so clean it up
#  - replace \\n with a new line: jq uses string value, replace it for real version
#  - remove " at starting of line: jq outputs wrapping string quotes, strip those out
########
function generate_aws_command() {
    local source_file="${DATA_SOURCE}"
jq 'map(
"  aws-vault exec \(.label)-\(.environment)-breakglass -- ./${binary_file} ~
    -day=|${day}|~
    -account_unit=|\(.billing_unit)|
") | join("\n")' ${source_file} \
    | sed 's/~/ \\/g' \
    | sed 's/|/"/g' \
    | sed 's/-null-/-/g' \
    | sed 's/\\n/\
/g' | sed 's/^"\(.*\)/\1/'

}

# generate string of the aws upload
function generate_upload_command() {
    echo '
read -p "upload to s3?: (Y|N) " up
if [[ "${up}" == [Yy] ]]; then
  cd ./data/
  ls -lh . | wc -l
  aws-vault exec ${profile} -- aws s3 cp --recursive . s3://${bucket}/${bucket_path} --sse AES256
  cd ../
fi'
}

# generate the complete bash file
function generate_file() {
    local file="${TARGET_FILE}"

    generate_file_head > ${file}
    generate_file_date_input >> ${file}
    generate_file_date_loop_start >> ${file}
    generate_aws_command >> ${file}
    generate_file_date_loop_end >> ${file}
    generate_upload_command >> ${file}

    chmod 0777 ${file}
    echo "✅ Generated cost fetching script"
    echo "  - [${file}]"
}

# move the binary we want to correct location
function move_binary() {
    local source="${BUILD_DIR}/reports/${BINARY}"
    local dest="${WORKING_DIR}"

    cp ${source} ${dest}
    echo "✅ Copied binary"
    echo "  - [${source}]"
    echo "  - [${dest}]"
}

function cleanup() {
    local release_dir="${RELEASE_DIR}"
    rm -Rf ${release_dir}
    echo "✅ Cleaned up release assets"
    echo "  - [${release_dir}]"
}

################################################
# MAIN FUNCTION
################################################
function main() {
    echo "------------------------------------------"
    download_realease
    generate_file
    move_binary
    cleanup
    echo "------------------------------------------"
    echo "Used these dates:"
    echo " - start_date: ${start_date}"
    echo " - end_date: ${end_date}"
    echo "------------------------------------------"
    if [[ "${OPEN_VSCODE}" == [Yy] ]]; then
        code ${TARGET_FILE}
    fi
}

main
