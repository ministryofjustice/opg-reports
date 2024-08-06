#!/usr/bin/env bash
# This script generates another script that contains a series of bash commands fetch
# monthly aws costs for all known accounts.
#
# Fetches the accounts.aws.json data from the latest opg-metadata release, and uses
# jq and sed to parse and convert the content into command to run
#

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/')
BUILD_ARCH="${OS}_${ARCH}"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BASE="${SCRIPT_DIR}/../builds/scripts"
BUILD_DIR="${SCRIPT_DIR}/../builds/${BUILD_ARCH}"
WORKING_DIR="${BASE}/local-aws-costs"

# prep a working directory
rm -Rf ${WORKING_DIR}
mkdir -p ${WORKING_DIR}
cd ${WORKING_DIR}

################################################
# DOWNLOAD AND ADJUST ACCOUNT METADATA FILE
################################################

# download the release of opg-metadata
mkdir -p ./metadata
cd ./metadata
gh release download --clobber --repo ministryofjustice/opg-metadata --pattern "*.tar.gz"
tar -xzf metadata.tar.gz
# use jq to remap all the values
echo "#!/usr/bin/env bash" > ./costs.sh
# add array of all months we want to record
echo 'months=("2023-06" "2023-07" "2023-08" "2023-09" "2023-10" "2023-11" "2023-12" "2024-01" "2024-02" "2024-03" "2024-04" "2024-05" "2024-06")' >> ./costs.sh
# add the for loop
echo 'for month in ${months[@]}; do' >> ./costs.sh
echo '  echo "month:${month}"' >> ./costs.sh
# echo 'month="2023-06"
# if [[ "${1}" != "" ]]; then
#     month="${1}"
# fi' >> ./costs.sh
########
# due to fun with quotes and escaping we use some subs with sed
# fixing it after jq. Details in order:
#  - replace ~ with \: so the command is neatly split over multiple lines
#  - replace | with ": so argument values have string encapsulation for spaces etc
#  - replace -null- with -: some accounts dont have an environment value, so clean it up
#  - replace \\n with a new line: jq uses string value, replace it for real version
#  - remove " at starting of line: jq outputs wrapping string quotes, strip those out
########
jq 'map(
"  aws-vault exec \(.label)-\(.environment)-breakglass -- ./aws_cost_monthly ~
    -month=|${month}|~
    -account_id=|\(.id)|~
    -account_name=|\(.name)|~
    -account_unit=|\(.billing_unit)|~
    -account_label=|\(.label)|~
    -account_environment=|\(.environment)|
"
) | join("\n")' ./accounts.aws.json \
    | sed "s/-month=|MONTH|/-month=|${month}|/g" \
    | sed 's/~/ \\/g' \
    | sed 's/|/"/g' \
    | sed 's/-null-/-/g' \
    | sed 's/\\n/\
/g' | sed 's/^"\(.*\)/\1/' >> costs.sh

# end the for loop
echo 'done' >> ./costs.sh

#####
# add the aws sync
#####
profile="shared-development-operator"
bucket="report-data-development"
bucket_path="aws/cost/monthly"
echo 'read -p "upload to s3?: (Y|N) " up' >> costs.sh
echo 'if [[ "${up}" == [Yy] ]]; then' >> costs.sh
echo "  cd ./data/" >> costs.sh
echo "  ls -lh . | wc -l" >> costs.sh
echo "  aws-vault exec ${profile} -- aws s3 cp --recursive . s3://${bucket}/${bucket_path} --sse AES256" >> costs.sh
echo "  cd ../" >> costs.sh
echo "  rm -Rf ./data/" >> costs.sh
echo 'fi' >> costs.sh
# mv the updated script up a level and clean up the directory
cp ./costs.sh ${BASE}/aws-costs-monthly.sh
chmod 0777 ${BASE}/aws-costs-monthly.sh
cd ../
rm -Rf ./metadata
################################################
# COPY OVER BINARY
################################################
cp ${BUILD_DIR}/reports/aws_cost_monthly ${BASE}


################################################
# cleanup
rm -Rf ${WORKING_DIR}


# echo out message
echo "------------------------------------------"
echo "✅ Download latest opg-metadata release and reformatted the account data"
echo "✅ Copied binary to same directory as script"
echo "Generated script here:"
echo "${BASE}/aws-costs-monthly.sh"
echo "------------------------------------------"
echo "Before use, you will need to:"
echo " - replace PROFILE with accurate aws profile values."
echo " - remove any non-aws account details"
echo " - update month to which ever reporting on"
echo "------------------------------------------"

code ${BASE}/aws-costs-monthly.sh
