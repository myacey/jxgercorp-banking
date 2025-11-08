.DEFAULT_GOAL=help

help:
	@echo "make up	- start docker containers"

up:
	docker compose --profile dev up --build -d
