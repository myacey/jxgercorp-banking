CREATE TYPE currency_enum AS ENUM ('RUB', 'USD', 'EUR');

CREATE TABLE "accounts" (
    "id" UUID PRIMARY KEY,
    "owner_username" VARCHAR(50) NOT NULL,
    "balance" BIGINT NOT NULL CHECK("balance" > 0),
    "currency" currency_enum NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT(NOW())
);

CREATE INDEX ON "accounts" ("owner_username");
CREATE INDEX ON "accounts" ("owner_username", "currency");

CREATE TABLE "transfers" (
    "id" UUID PRIMARY KEY,
    "from_account_id" UUID NOT NULL REFERENCES "accounts" ("id"),
    "to_account_id" UUID NOT NULL REFERENCES "accounts" ("id"),
    "amount" BIGINT NOT NULL CHECK("amount" > 0),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT(NOW())
);

CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
