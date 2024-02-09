run:
	APP_PORT="9999" DATABASE_URL="postgres://admin:123@127.0.0.1:5432/rinha?pool_max_conns=10" go run ./src/*.go

create-image:
	docker build --tag gabrielucelli .

compose-up:
	docker compose up

compose-down:
	docker compose down