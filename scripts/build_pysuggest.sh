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
else
    echo "Directory does not exist, build failed!"
fi
