#!/bin/sh
if [ ! -d $1 ];
then 
    echo "Checking change failed, '$1' is invalid dir path" >&2
    exit 1
fi

echo "Checking GO mod change: $1"

SELF_GIT_STATUS=$(git status --porcelain $1)

if [ -n "$SELF_GIT_STATUS" ];
then
    echo '{ "change": "self" }'
    exit 0
fi

LOCAL_DEPENDENCIES=$(grep -r '"blinders/.*"' $1/*.go $1/**/*.go  | sed -n 's/.*\(blinders.*\)"/\1/p' | sort | uniq)

while IFS= read -r module_name; 
do
    echo "Checking dependency: $module_name"

    module_string="module $module_name"
    module_path=$(grep -l "$module_string" **/go.mod **/**/go.mod | head -n 1)

    # ignore un-resolved module
    if [ -z $module_path ];
    then 
        echo "Not found path for mod: $module_name, continue"
        continue
    fi
    
    DEPENDENCY_GIT_STATUS=$(git status --porcelain $module_path)

    if [ -n "$DEPENDENCY_GIT_STATUS" ];
    then
        echo "{ \"change\": \"dependency\", \"path\": \"$module_path\" }"
        exit 0
    fi
done <<< "$LOCAL_DEPENDENCIES"

echo '{ "change": "unchange" }'