.PHONY: run e2e

run:
	COMPOSE_BAKE=true docker compose up --build

e2e:
	COMPOSE_BAKE=true docker compose -f docker-compose.e2e.yml up -d --build
	pytest -q tests/test_teams.py tests/test_pull_requests.py tests/test_users.py || true
	docker compose -f docker-compose.e2e.yml down -v

down:
	docker compose -f docker-compose.yml down -v

lint:
	golangci-lint run ./...
