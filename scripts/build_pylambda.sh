# !/bin/bash

if [[ "$1" != "dev" && "$1" != "staging" && "$1" != "prod" ]]; 
then
    echo "Usage: $0 with one of 'dev|staging|prod'"
    exit 1
fi

cd functions/dictionary
poetry build
poetry run pip install --upgrade -t bundle dist/*.whl
cd bundle ; zip -r ../lambda_bundle.zip . -x '*.pyc'
cd ..
rm -rf ./dist ./bundle

cp ./lambda_bundle.zip ../../dist/dictionary-$1.zip