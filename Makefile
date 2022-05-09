run-tests:
		docker volume prune -f && \
		docker compose -f ./docker-compose.test.yaml build && \
		docker compose -f ./docker-compose.test.yaml \
		run test-rider-service gotest -v -p=1 ./... && \
		docker compose -f ./docker-compose.test.yaml down

golangci:
	golangci-lint run