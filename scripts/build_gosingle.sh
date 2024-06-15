if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ./dist/gosuggest-$1/bootstrap ./services/suggest/lambda
echo "build gosuggest lambda function completed"
cp ./firebase.admin.$1.json ./dist/gosuggest-$1/firebase.admin.json
echo "copied firebase.admin.json to gosuggest"
cd ./dist/gosuggest-$1
zip -r ../gosuggest-$1.zip .
cd ../..