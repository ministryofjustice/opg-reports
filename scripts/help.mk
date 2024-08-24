#!/bin/bash
RULE_COLOR="$(tput setaf 6)"
SECTION_COLOR="$(tput setaf 3)"
VARIABLE_COLOR="$(tput setaf 2)"
VALUE_COLOR="$(tput setaf 1)"
CLEAR_STYLE="$(tput sgr0)"
TARGET_STYLED_HELP_NAME="${RULE_COLOR}TARGET${CLEAR_STYLE}"
ARGUMENTS_HELP_NAME="${VARIABLE_COLOR}ARGUMENT${CLEAR_STYLE}=${VALUE_COLOR}VALUE${CLEAR_STYLE}"

echo "${SECTION_COLOR}Targets:${CLEAR_STYLE}"
echo "    ${RULE_COLOR}help${CLEAR_STYLE}"
echo "        Get help for commands in this folder"
echo ""

TARGET_REGEX="^[a-zA-Z0-9%_\/%-]+:"
SECTION_REGEX="^##\s*@section\s*(.*)$"
DOCBLOCK_REGEX="^##\s*(.*)$"
PARAM_REGEX="@param\s+([a-zA-Z_]+)(=([^\s]+))?\s*(.*$)?"

COMMENT=""
PARAMS=""
PARAMS_DOC=""
cat $FILE | while read line
do
    # do something with $line here
    if [[ ! -z $line ]]
    then
        if [[ $line =~ $SECTION_REGEX ]]
        then
            SECTION_NAME=$(echo $line | sed -e "s/^##\s*@section\s*\(.*\)$/\1/g")
            echo "$SECTION_COLOR$SECTION_NAME$CLEAR_STYLE:"
        elif [[ $line =~ $TARGET_REGEX ]]
        then
            # if there is no comment for this target, we don't display it in the docs to keep private targets hidden
            if [[ ! -z $COMMENT ]]
            then
                TARGET=$(echo $line | sed -e "s/^\([a-zA-Z0-9%_\/%-]\+\):.*/\1/g")
                echo "    $RULE_COLOR$TARGET$CLEAR_STYLE $PARAMS"
                echo -e "$COMMENT"
                if [[ ! -z $PARAMS_DOC ]]
                then
                    echo "        Params:"
                    echo -e "$PARAMS_DOC"
                fi
            fi
            COMMENT=""
            PARAMS=""
            PARAMS_DOC=""
        elif [[ $line =~ $PARAM_REGEX ]]
        then
            PARAM=$(echo $line | sed -e "s/##\s*@param\s\+\([a-zA-Z_]\+\)\(=\([^[:space:]]\+\)\)\?\s*\(.*\)\?$/${VARIABLE_COLOR}\1${CLEAR_STYLE}=${VALUE_COLOR}\3${CLEAR_STYLE}/g")
            PARAM_DOC=$(echo $line | sed -e "s/##\s*@param\s\+\([a-zA-Z_]\+\)\(=\([^[:space:]]\+\)\)\?\s*\(.*\)\?$/- \1 (ex: \3) \4/g")
            PARAMS="${PARAMS}${PARAM} "
            PARAMS_DOC="${PARAMS_DOC}         ${PARAM_DOC}\n"
        elif [[ $line =~ $DOCBLOCK_REGEX ]]
        then
            # echo "doc : $line"
            # echo $line | sed -e "s/^##\s*\(.*\)$/\1/g"
            LINE_CLEANED=$(echo $line | sed -e "s/^##\s*\(.*\)$/\1/g")
            COMMENT="${COMMENT}        $LINE_CLEANED\n"
        fi
    fi
done
