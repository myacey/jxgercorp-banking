CREATE TABLE currencies (
  code TEXT PRIMARY KEY,
  symbol TEXT NOT NULL,
  precision INT NOT NULL
);

INSERT INTO currencies VALUES
    ('RUB', '₽', 2),
    ('USD', '$', 2),
    ('EUR', '€', 2);

ALTER TABLE accounts
ADD COLUMN currency_code TEXT;

UPDATE accounts
SET currency_code = currency::text;

ALTER TABLE accounts
ALTER COLUMN currency_code SET NOT NULL;

ALTER TABLE accounts
ADD CONSTRAINT fk_accounts_currency
FOREIGN KEY (currency_code) REFERENCES currencies(code);

ALTER TABLE accounts
DROP COLUMN currency;

ALTER TABLE transfers
ADD COLUMN currency_code TEXT,
ADD COLUMN from_account_username TEXT,
ADD COLUMN to_account_username TEXT;

UPDATE transfers
SET
    currency_code = accounts.currency_code,
    from_account_username = accounts.owner_username
FROM accounts
WHERE transfers.from_account_id = accounts.id;

UPDATE transfers
SET
    currency_code = accounts.currency_code,
    to_account_username = accounts.owner_username
FROM accounts
WHERE transfers.to_account_id = accounts.id;

ALTER TABLE transfers
ALTER COLUMN currency_code SET NOT NULL,
ALTER COLUMN from_account_username SET NOT NULL,
ALTER COLUMN to_account_username SET NOT NULL;

ALTER TABLE transfers
ADD CONSTRAINT fk_transfers_currency
FOREIGN KEY (currency_code) REFERENCES currencies(code);

DROP TYPE currency_enum;
