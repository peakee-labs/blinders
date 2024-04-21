<h1 align="center">Blinders</h1>
<p align="center">Monorepo, microservices backend for a language learning platform - Peakee.</p>
<p align="center">Golang, Python, Langchain, AWS, Terraform, Firebase, MongoDB</p>

## System Architecture

<img width="1081" alt="image" src="https://github.com/dev-zenonian/blinders/assets/104194494/91616345-53d9-4675-9a0a-d2e8b7646d0c">

## Resources

For local development, we use `.env` as default for running services locally and for testing. Have a look at `.env.example`. Also, we use `firebase.admin.json` for firebase authentication integration in local.

For deployment with stages, we use `.env.dev`, `.env.staging`, `.env.prod` for corresponding stages `dev`, `staging`, `prod`. With firebase, we use `firebase.admin.dev.json`, `firebase.admin.staging.json` and `firebase.admin.prod.json`

## Go development setup

Require Go version >= 1.21

Using [golangci-lint](https://golangci-lint.run/) to manage all linter, formatter and setup ci, detail configs in `golangci.yml`. You should config `golangci-lint` in your code editor to pass all the linters

### Live-reloading

Install [air](https://github.com/cosmtrek/air) for live-reloading when working with Go

```
go install github.com/cosmtrek/air@latest
```

## Python development setup

### Tools

- Code formatter: [Black](https://github.com/psf/black)
- Code linter: [Flake8](https://flake8.pycqa.org/en/latest/user/index.html) [isort](https://github.com/PyCQA/isort), and [pylint](https://pypi.org/project/pylint/) for just checking public artifacts are documented or not
- Type checking: [Pyright](https://github.com/microsoft/pyright#static-type-checker-for-python)

### Setup steps

Use Python 3.10 as the base version of Python, recommend to use a local Python environment using [conda](https://www.anaconda.com/)

```shell
conda create --prefix ./.venv/ python==3.10
conda activate ./.venv
```

We're using [poetry](https://python-poetry.org/) package manager because of rich dependencies management features

```shell
pip install poetry && poetry install
```

If not using `poetry`

```shell
pip install -e .
```

## CLI tools

Install the CLI to go packages

```
make setup-cli
```

Need to setup `.env`, use `.env.production` and `.env.development`. See example in `env.example`

### Usage

Run help command for more details

```
blinders --help
```

To work with `auth` commands:

```
# get jwt of a user
blinders auth load-user --uid <user_uid>
```

```
# generate wscat command to connect as a client
# default endpoint ws.peakee.co/v1
blinders auth gen-wscat --endpoint <endpoint> --uid <user_uid>
```

## Local development

Run development docker compose to prepare development environment

```
make dev-container
```

Run REST api with Air

```
make rest
```

Run Embedder service

```
make embedder
# or
poetry run embedder_service
```

## Deployment

We have 3 separated environments which are `dev`, `staging` and `prod`. Each deployment environment have a deployment state.

Also we have a `shared` state for managing `route53_certificate` and deploying some shared ec2 instances for `dev` and `staging` environments

Use `.env.dev|staging|prod`, `firebase.admin.dev|staging|prod.json` at root corresponding to `dev|staging|prod`

Pre-build functions

```
sh scripts/build_golambda.sh dev|staging|prod
sh scripts/build_pylambda.sh dev|staging|prod
sh scripts/build_pysuggest.sh dev|staging|prod
```

At the first time or having any update to shared state

```
cd infra/shared && terraform plan
terraform apply
```

For a specific environment

```
cd infra/dev|staging|prod && terraform plan
terraform apply
```
