if [ ! -d $1 ];
then 
    echo "$1 is invalid dir path"
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
    echo "Checking $module_name"
    grep -r "module $module_name"
    
    DEPENDENCY_GIT_STATUS=$(git status --porcelain $module_name)
    echo $DEPENDENCY_GIT_STATUS
    if [ -n "$DEPENDENCY_GIT_STATUS" ];
    then
        echo '{ "change": "dependency", "path": "$module_name" }'
        exit 0
    fi
done <<< "$LOCAL_DEPENDENCIES"

echo '{ "change": "unchange" }'