DOCKER_BUILDKIT=1 docker build -t blinders-pysuggest-lambda . -f Dockerfile.build-pylambda -o ./functions/suggest/lambda_bundle

cd functions/suggest/lambda_bundle
zip -r ../lambda_bundle.zip . -x '*.pyc'