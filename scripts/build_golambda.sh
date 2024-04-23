# !/bin/bash
# TODO: separate env builds

# GOOS is the target operating system (linux, darwin, etc.)
# GOARCH is the target architecture (386, amd64, etc.)
# CGO_ENABLED=0 disables cgo (linking of C libraries)
# GOFLAGS=-trimpath removes debug info from the binary
# -mod=readonly disallows updating go.mod and go.sum
# -ldflags='-s -w' strips symbol table and debug info from the binary

# need to zip at the target directory with "."

if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

rm -rf dist/connect*$1 dist/translate*$1 dist/authorizer*$1 \
    dist/explore*$1 dist/disconnect*$1 dist/wschat*$1$1 \
    dist/rest*$1 dist/notification*$1 dist/ws_authorizer*$1 \
    dist/collecting*$1 

echo "cleaned previous build artifacts"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/translate-$1/bootstrap ./functions/translate
echo "build translate lambda function completed"
cp ./firebase.admin.$1.json ./dist/translate-$1/firebase.admin.json
echo "copied firebase.admin.json to translate"
cd ./dist/translate-$1
zip -r ../translate-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/connect-$1/bootstrap ./functions/websocket/connect
echo "build connect lambda function completed"
cd ./dist/connect-$1
zip -r ../connect-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/authorizer-$1/bootstrap ./functions/websocket/authorizer
echo "build authorizer lambda function completed"
cp ./firebase.admin.$1.json ./dist/authorizer-$1/firebase.admin.json
echo "copied firebase.admin.json to authorizer"
cd ./dist/authorizer-$1
zip -r ../ws_authorizer-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/disconnect-$1/bootstrap ./functions/websocket/disconnect
echo "build disconnect lambda function completed"
cd ./dist/disconnect-$1
zip -r ../disconnect-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/wschat-$1/bootstrap ./functions/websocket/chat
echo "build websocket chat lambda function completed"
cd ./dist/wschat-$1
zip -r ../wschat-$1.zip .
cd ../..

# migrate to arm64 for better price-performance
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/rest-$1/bootstrap ./functions/rest
echo "build rest api lambda function completed"
cp ./firebase.admin.$1.json ./dist/rest-$1/firebase.admin.json
echo "copied firebase.admin.json to rest api"
cd ./dist/rest-$1
zip -r ../rest-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/notification-$1/bootstrap ./functions/websocket/notification
echo "build notification function completed"
cd ./dist/notification-$1
zip -r ../notification-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/explore-$1/bootstrap ./functions/explore
echo "build explore lambda function completed"
cp ./firebase.admin.$1.json ./dist/explore-$1/firebase.admin.json
echo "copied firebase.admin.json to explore api"
cd ./dist/explore-$1
zip -r ../explore-$1.zip .
cd ../..

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/collecting-$1/bootstrap ./functions/collecting
echo "build collecting lambda function completed"
cp ./firebase.admin.$1.json ./dist/collecting-$1/firebase.admin.json
echo "copied firebase.admin.json to collecting api"
cd ./dist/collecting-$1
zip -r ../collecting-$1.zip .
cd ../..