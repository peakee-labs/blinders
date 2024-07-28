#!/bin/sh

if [[ "$2" != "dev" && "$2" != "staging" && "$2" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

if [ ! -d "$3" ]; then
    echo "Usage: $3 must be a dir"
    exit 1
fi

if [[ "$4" != "firebase" && "$4" != "none" ]]; 
then
    echo "Usage: $4 with one of 'firebase|none', firebase value for copying firebase.admin.json"
    exit 1
fi

if [ -d "./dist/$1-$2" ];
then
    rm -rf ./dist/$1-$2
    echo "Removed previous $1-$2 build"
fi

echo "Building Go Module $1-$2"
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/$1-$2/bootstrap $3
echo "Completed building Go Module $1-$2"

if [[ "$4" == "firebase" ]];
then
    cp ./firebase.admin.$2.json ./dist/gosuggest-$2/firebase.admin.json
    echo "Copied firebase.admin.json to $1-$2"
fi

cd ./dist/$1-$2
zip -r ./bundle.zip .
cd ../..

echo "Completed to ./dist/$1-$2 on $3 with $4"
