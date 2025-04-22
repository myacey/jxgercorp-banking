include .env
export

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# postgres:
# 	docker run --name postgres17 -p $(POSTGRES_PORT):$(POSTGRES_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres:17-alpine

# createdb:
# 	docker exec -it postgres17 createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

# dropdb:
# 	docker exec -it postgres17 dropdb -U root $(POSTGRES_DB)

# migrateup:
# 	migrate -path db/migration -database "${DB_URL}" -verbose up

# migratedown:
# 	migrate -path db/migration -database "${DB_URL}" -verbose down

sqlc:
	sqlc generate

# redis:
# 	docker run --name redis -p 6379:6379 -d redis:7-alpine

new_migration:
	migrate create -ext sql -dir $(dir) -seq $(name)
# make new_migration name=<some name>

.PHONY: postgres createdb dropdb migrateup migratedown sqlc up redis new_migration
