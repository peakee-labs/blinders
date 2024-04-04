# Blinders

Monorepo, a microservices backend project named Blinders for a language learning platform. Almost use Golang, Python, and Terraform for provisioning, all the services are hosted on AWS.

## System Architecture

<img width="1081" alt="image" src="https://github.com/dev-zenonian/blinders/assets/104194494/91616345-53d9-4675-9a0a-d2e8b7646d0c">

## Resources

Prepare env file `.env.development`, please check from `.env.example` for more detail

Prepare firebase admin credentials file, create file `firebase.admin.development.json`

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

# Deployment

Prepare environment file `.env`

Pre-build suggest function

```
sh ./scripts/build_pysuggest.sh
```

Deploy

```
cd infra && terraform apply
```
