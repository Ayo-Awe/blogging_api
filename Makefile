include .env

dev:
	air

migrate-up:
	migrate -path migrations -verbose -database ${DB_URL} up

migrate-down:
	migrate -path migrations -verbose -database ${DB_URL} down

migrate-force:
	migrate -path migrations -database ${DB_URL} -verbose force ${version}

new-migration:
	migrate create -dir migrations -seq -ext .sql ${name}

migrate-drop:
	migrate -path migrations -database ${DB_URL} -verbose drop

.PHONY: migrate-up migrate-down new-migration migrate-force dev
