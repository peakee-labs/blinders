# !/bin/sh

PY_FILES=$(find . -type f -name "*.py" -not -path "./.venv/*" -not -path "**/lambda_bundle/**")

echo "-> Lint by isort"
isort --check-only $PY_FILES

echo "-> Lint by black"
black --check $PY_FILES

echo "-> Lint by flake8"
flake8 $PY_FILES

echo "-> Lint by pyright"
pyright $PY_FILES
