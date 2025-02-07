#!/bin/sh
sleep 20
migrate -path /migrations -database "postgresql://root:secret@postgres:5432/jxger_bank?sslmode=disable" up
