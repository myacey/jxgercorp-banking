sqlc:
	sqlc generate -f ./db/sqlc.yaml


new_migration:
	migrate create -ext sql -dir ./db/migrations -seq $(name)

