<h1 align="center">Blinders</h1>
<p align="center">Monorepo, microservices backend for a language learning platform - Peakee.</p>
<p align="center">Golang, Python, Langchain, AWS, Terraform, Firebase, MongoDB</p>

## System Architecture

<img width="1081" alt="System architecture" src="https://github.com/zenonian-labs/blinders/assets/104194494/14e735a6-19b0-4dc0-a22f-c49001b753cc">

## Resources

For local development, we use `.env` as the default for running services locally and for testing. Have a look at `.env.example`. Also, we use `firebase.admin.json` for firebase authentication integration in local.

For deployment with stages, we use `.env.dev`, `.env.staging`, `.env.prod` for corresponding stages `dev`, `staging`, `prod`. With firebase, we use `firebase.admin.dev.json`, `firebase.admin.staging.json`, and `firebase.admin.prod.json` respectively.

Just a convention, we try to use `dev stage` nearly the same as `local development`, e.g. firebase admin, API keys,...

## Go development setup

Require Go version >= 1.21

Using [golangci-lint](https://golangci-lint.run/) to manage all linter, formatter, and setup ci, detail configs in `golangci.yml`. You should config `golangci-lint` in your code editor to pass all the linters

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

Need to setup `.env`, use `.env.production` and `.env.development`. See the example in `env.example`

### Usage

Run the help command for more details

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

Run development docker-compose to prepare the development environment

```
make dev-container
```

Run REST API with Air

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

We have 3 separate environments:  `dev`, `staging`, and `prod`. Each deployment environment has a deployment state.

Also, we have a `shared` state for managing `route53_certificate` and deploying some shared ec2 instances for `dev` and `staging` environments

Use `.env.dev|staging|prod`, `firebase.admin.dev|staging|prod.json` at root corresponding to `dev|staging|prod`

We use `s3` backend to store the deployment state, you need to set up the config file `infra/backend.conf`, see backend.conf.example. Remember to pre-create bucket and dynamodb table following [setup s3 backend](https://developer.hashicorp.com/terraform/language/settings/backends/s3). For each state, init terraform at the first time:

```
cd infra/shared|dev|staging|prod
terraform init -backend-config=../backend.conf
```

Pre-build functions

```
sh scripts/build_golambda.sh dev|staging|prod
sh scripts/build_pylambda.sh dev|staging|prod
sh scripts/build_pysuggest.sh dev|staging|prod

# or build all functions
sh scripts/build_all.sh dev|staging|prod
```

At the first time or having any update to the shared state

```
cd infra/shared && terraform plan
terraform apply

# or with a profile
terraform apply -var="profile=..."
```

For a specific environment

```
cd infra/dev|staging|prod && terraform plan
terraform apply

# or with a profile
terraform apply -var="profile=..."
```

Otherwise, we have `trydev` stage which is used for deployment testing with the same env as `dev` stage, use this stage to test your new updating deployment before creating PR. Build step and resource files are the same as `dev`

```
cd infra/trydev && terraform plan
terraform apply

# or with a profile
terraform apply -var="profile=..."
```
