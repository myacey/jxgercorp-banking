CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar UNIQUE NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "balance" bigserial NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
    "id" bigserial PRIMARY KEY,
    "from_user" varchar NOT NULL,
    "to_user" varchar NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_user") REFERENCES "users" ("username");
ALTER TABLE "transactions" ADD FOREIGN KEY ("to_user") REFERENCES "users" ("username");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "transactions" ("from_user");

CREATE INDEX ON "transactions" ("to_user");
