FROM migrate/migrate:latest

COPY services/user/db/migrations        /migrations/user/
COPY services/transfer/db/migrations    /migrations/transfer/

ENTRYPOINT ["/usr/bin/env"]
