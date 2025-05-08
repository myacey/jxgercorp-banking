CREATE TYPE user_status_enum AS ENUM ('pending', 'active', 'banned');

CREATE TABLE "users" (
    "id" UUID PRIMARY KEY,
    "username" varchar(50) UNIQUE NOT NULL,
    "email" varchar(50) UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "status" user_status_enum NOT NULL DEFAULT('pending')
);

CREATE INDEX ON "users" ("username");