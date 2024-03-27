# !/bin/bash

# GOOS is the target operating system (linux, darwin, etc.)
# GOARCH is the target architecture (386, amd64, etc.)
# CGO_ENABLED=0 disables cgo (linking of C libraries)
# GOFLAGS=-trimpath removes debug info from the binary
# -mod=readonly disallows updating go.mod and go.sum
# -ldflags='-s -w' strips symbol table and debug info from the binary

rm -rf dist
echo "cleaned dist directory"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/translate/bootstrap ./functions/translate
echo "build translate lambda function completed"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/connect/bootstrap ./functions/websocket/connect
echo "build connect lambda function completed"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/authorizer/bootstrap ./functions/websocket/authorizer
echo "build authorizer lambda function completed"
cp ./firebase.admin.json ./dist/authorizer
echo "copied firebase.admin.json to authorizer"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/disconnect/bootstrap ./functions/websocket/disconnect
echo "build disconnect lambda function completed"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/wschat/bootstrap ./functions/websocket/chat
echo "build websocket chat lambda function completed"

# migrate to arm64 for better price-performance
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/rest/bootstrap ./functions/rest
echo "build rest api lambda function completed"
cp ./firebase.admin.json ./dist/rest
echo "copied firebase.admin.json to rest api"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/notification/bootstrap ./functions/websocket/notification
echo "build notification function completed"

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ./dist/explore/bootstrap ./functions/explore
echo "build explore lambda function completed"
cp ./firebase.admin.json ./dist/explore
echo "copied firebase.admin.json to explore api"