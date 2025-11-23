lint:
	golangci-lint run --fix

gen:
	go generate ./...

up:
	docker-compose up

build:
	docker-compose up --build

test:
	IS_TEST=true \
	docker-compose --file docker-compose.yml --file docker-compose-test.yml up --build --exit-code-from test --abort-on-container-exit

	docker logs -f test
