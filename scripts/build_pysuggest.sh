if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

DOCKER_BUILDKIT=1 docker build -t blinders-pysuggest-lambda . -f Dockerfile.build-pylambda -o ./functions/suggest/lambda_bundle

cd functions/suggest/lambda_bundle
zip -r ../lambda_bundle.zip . -x '*.pyc'
cd ..
cp ./lambda_bundle.zip ../../dist/suggest-$1.zip