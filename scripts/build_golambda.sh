# !/bin/bash
# TODO: separate env builds

# GOOS is the target operating system (linux, darwin, etc.)
# GOARCH is the target architecture (386, amd64, etc.)
# CGO_ENABLED=0 disables cgo (linking of C libraries)
# GOFLAGS=-trimpath removes debug info from the binary
# -mod=readonly disallows updating go.mod and go.sum
# -ldflags='-s -w' strips symbol table and debug info from the binary
if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

# rm -rf dist
# echo "cleaned dist directory"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/translate-$1/bootstrap ./functions/translate
echo "build translate lambda function completed"
zip -r ./dist/translate-$1.zip ./dist/translate-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/connect-$1/bootstrap ./functions/websocket/connect
echo "build connect lambda function completed"
zip -r ./dist/connect-$1.zip ./dist/connect-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/authorizer-$1/bootstrap ./functions/websocket/authorizer
echo "build authorizer lambda function completed"
cp ./firebase.admin.$1.json ./dist/authorizer-$1/firebase.admin.json
echo "copied firebase.admin.json to authorizer"
zip -r ./dist/ws_authorizer-$1.zip ./dist/authorizer-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/disconnect-$1/bootstrap ./functions/websocket/disconnect
echo "build disconnect lambda function completed"
zip -r ./dist/disconnect-$1.zip ./dist/disconnect-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/wschat-$1/bootstrap ./functions/websocket/chat
echo "build websocket chat lambda function completed"
zip -r ./dist/wschat-$1.zip ./dist/wschat-$1

# migrate to arm64 for better price-performance
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/rest-$1/bootstrap ./functions/rest
echo "build rest api lambda function completed"
cp ./firebase.admin.$1.json ./dist/rest-$1/firebase.admin.json
echo "copied firebase.admin.json to rest api"
zip -r ./dist/rest-$1.zip ./dist/rest-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/notification-$1/bootstrap ./functions/websocket/notification
echo "build notification function completed"
zip -r ./dist/notification-$1.zip ./dist/notification-$1

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/explore-$1/bootstrap ./functions/explore
echo "build explore lambda function completed"
cp ./firebase.admin.$1.json ./dist/explore/firebase.admin.json
echo "copied firebase.admin.json to explore api"
zip -r ./dist/explore-$1.zip ./dist/explore-$1