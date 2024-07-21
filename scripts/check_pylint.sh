# !/bin/bash

echo "-> Lint by isort"
isort --check-only $(git ls-files "**/*.py")

echo "-> Lint by black"
black --check $(git ls-files "**/*.py")

echo "-> Lint by pylint"
pylint $(git ls-files --full-name -- "**/*.py" | xargs -n 1 dirname | uniq)

echo "-> Lint by flake8"
flake8 $(git ls-files "**/*.py")

echo "-> Lint by pyright"
pyright $(git ls-files "**/*.py")