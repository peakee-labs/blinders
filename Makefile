golint:
	go list -f '{{.Dir}}/...' -m | xargs -n 1 sh -c 'golangci-lint run $$0 || exit -1'

gotest:
	go list -f '{{.Dir}}/...' -m | xargs -n 1 sh -c 'go test $$0 || exit -1'

dev-container:
	docker compose -f dev.docker-compose.yml up -d

dev-container-down:
	docker compose -f dev.docker-compose.yml down

setup-cli:
	go install ./cli/blinders.go

rest:
	air -c ./tools/.air.rest.toml

embedder:
	poetry run embedder_service