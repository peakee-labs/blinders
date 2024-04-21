
if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi


rm -rf dist/suggest*$1

suggest_bundle="functions/suggest/lambda_bundle"
suggest_bundle_zip="functions/suggest/lambda_bundle.zip"

if [ -d "$suggest_bundle" ]; then
    echo "Cleaning build directory..."
    rm -rf functions/suggest/lambda_bundle
fi

if [ -f "$suggest_bundle_zip" ]; then
    echo "Cleaning zip file..."
    rm functions/suggest/lambda_bundle.zip
fi

DOCKER_BUILDKIT=1 docker build -t blinders-pysuggest-lambda . -f Dockerfile.build-pylambda -o ./functions/suggest/lambda_bundle

if [ -d "$suggest_bundle" ]; then
    echo "Suggest bundle exists"
    cd $suggest_bundle
    zip -r ../lambda_bundle.zip . -x '*.pyc'
    cd ..
    cp ./lambda_bundle.zip ../../dist/suggest-$1.zip
else
    echo "Directory does not exist, build failed!"
fi

cd ../..