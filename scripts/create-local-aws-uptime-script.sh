#!/usr/bin/env bash
# This script generates another script that contains a series of bash commands to
# fetch historical uptime data from our accounts
#
# Fetches the accounts.aws.uptime.json data from the latest opg-metadata release,
# and uses jq and sed to parse and convert the content into command to run
#
# Will ask for user input for start and end days

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/')
BUILD_ARCH="${OS}_${ARCH}"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BASE="${SCRIPT_DIR}/../builds/scripts"
BUILD_DIR="${SCRIPT_DIR}/../builds/${BUILD_ARCH}"
WORKING_DIR="${BASE}/local-aws-uptime"

# prep a working directory
rm -Rf ${WORKING_DIR}
mkdir -p ${WORKING_DIR}
cd ${WORKING_DIR}

################################################
# DOWNLOAD AND ADJUST UPTIME METADATA FILE
################################################

# download the release of opg-metadata
mkdir -p ./metadata
cd ./metadata
gh release download --clobber --repo ministryofjustice/opg-metadata --pattern "*.tar.gz"
tar -xzf metadata.tar.gz
# loop between two dates
echo '#!/usr/bin/env bash
read -p "Start date [2024-07-22]: " start
start="${start:-2024-07-22}"
read -p "End date [2024-08-06]: " end
end="${end:-2024-08-06}"
d="${start}"
while [ "$d" != "${end}" ]; do
  echo "day: ${d}"
  d=$(date -j -v +1d -f "%Y-%m-%d" $d +%Y-%m-%d)
' > ./uptime.sh



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
"  aws-vault exec \(.label)-\(.environment)-breakglass -- ./aws_uptime_daily ~
    -day=|${d}|~
    -account_unit=|\(.billing_unit)|
"
) | join("\n")' ./accounts.aws.uptime.json \
    | sed 's/~/ \\/g' \
    | sed 's/|/"/g' \
    | sed 's/-null-/-/g' \
    | sed 's/\\n/\
/g' | sed 's/^"\(.*\)/\1/' >> ./uptime.sh

# end the for loop
echo 'done' >> ./uptime.sh

# #####
# # add the aws sync
# #####
echo '
profile="shared-development-operator"
bucket="report-data-development"
bucket_path="aws/uptime/daily"
' >> ./uptime.sh
echo '
read -p "upload to s3?: (Y|N) " up
if [[ "${up}" == [Yy] ]]; then
  cd ./data/
  ls -lh . | wc -l
  aws-vault exec ${profile} -- aws s3 cp --recursive . s3://${bucket}/${bucket_path} --sse AES256
  cd ../
  rm -Rf ./data/
fi' >> ./uptime.sh
# # mv the updated script up a level and clean up the directory
cp ./uptime.sh ${BASE}/aws-uptime-daily.sh
chmod 0777 ${BASE}/aws-uptime-daily.sh
cd ../
rm -Rf ./metadata
# ################################################
# # COPY OVER BINARY
# ################################################
cp ${BUILD_DIR}/reports/aws_uptime_daily ${BASE}


# ################################################
# # cleanup
rm -Rf ${WORKING_DIR}


# # echo out message
echo "------------------------------------------"
echo "✅ Download latest opg-metadata release and reformatted the account data"
echo "✅ Copied binary to same directory as script"
echo "Generated script here:"
echo "${BASE}/aws-uptime-daily.sh"
echo "------------------------------------------"

code ${BASE}/aws-uptime-daily.sh
