.DEFAULT_GOAL = help

help:
	@echo "make sqlc			- generate sqlc files"
	@echo "make new_migration <name>	- create new migration files (up and down)"

sqlc:
	sqlc generate -f ./db/sqlc.yaml

new_migration:
	migrate create -ext sql -dir ./db/migrations -seq $(name)
