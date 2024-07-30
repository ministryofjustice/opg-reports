#!/usr/bin/env bash
# This script fetches the latest release of the opg-metadata, uses
# the xml version of account data to form the basis of the ./aws-costs.sh
# file
# After it has run, it you will need to edit the file and change PROFILE
# with correct profile values for the local host

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
cp ./accounts.xml ./costs.sh
# run string substitues on the xml tags to become arguments for the command
# map the account id
sed -i'' -e 's#<id>#-account_id "#g' ./costs.sh
sed -i'' -e 's#</id>#" \\#g' ./costs.sh
# map the account name
sed -i'' -e 's#<name>#-account_name "#g' ./costs.sh
sed -i'' -e 's#</name>#" \\#g' ./costs.sh
# map the unit
sed -i'' -e 's#<billing-unit>#-account_unit "#g' ./costs.sh
sed -i'' -e 's#</billing-unit>#" \\#g' ./costs.sh
# map the label
sed -i'' -e 's#<label>#-account_label "#g' ./costs.sh
sed -i'' -e 's#</label>#" \\#g' ./costs.sh
# map the environment
sed -i'' -e 's#<environment>#-account_environment "#g' ./costs.sh
sed -i'' -e 's#</environment>#"#g' ./costs.sh
# remove the type
sed -i'' -e 's#<type>aws</type>##g' ./costs.sh
# setup the command
sed -i'' -e 's#<account>#aws-vault exec PROFILE -- ./aws_cost_monthly -month "2024-05" \\#g' ./costs.sh
# remove closing tags
sed -i'' -e 's#</account>##g' ./costs.sh
sed -i'' -e 's@<accounts>@#!/usr/bin/env bash@g' ./costs.sh
sed -i'' -e 's#</accounts>##g' ./costs.sh
# remove tabs / spaces at start of the line
sed -i'' -e "s/^[ \t]*//" ./costs.sh
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
echo "✅ Download latest opg-metadata release and reformatted the accounts.xml file"
echo "✅ Copied binary to same directory as script"
echo "Generated script here:"
echo "${BASE}/aws-costs-monthly.sh"
echo "------------------------------------------"
echo "Before use, you will need to:"
echo " - replace PROFILE with accurate aws profile values."
echo " - remove any non-aws account details"
echo " - update month to which ever reporting on"
echo "------------------------------------------"

read -p "Edit file? (Y|N) " edit_file

if [[ "${edit_file}" == [yY] ]]; then
    code ${BASE}/aws-costs-monthly.sh
fi
