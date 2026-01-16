CREATE TYPE currency_enum AS ENUM ('RUB', 'USD', 'EUR');

ALTER TABLE accounts
ADD COLUMN currency currency_enum;

UPDATE accounts
SET currency = currency_code::currency_enum;

ALTER TABLE accounts
DROP CONSTRAINT fk_accounts_currency;

ALTER TABLE accounts
DROP COLUMN currency_code;

ALTER TABLE transfers
DROP CONSTRAINT fk_transfers_currency;

ALTER TABLE transfers
DROP COLUMN currency_code;

ALTER TABLE transfers
DROP COLUMN owner_username;

DROP TABLE currencies;
