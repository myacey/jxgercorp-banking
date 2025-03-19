BEGIN;

CREATE TYPE user_status AS ENUM ('pending', 'active', 'banned');

ALTER TABLE users
ADD COLUMN status user_status DEFAULT 'pending' NOT NULL;

UPDATE users
SET status = CASE
    WHEN pending THEN 'pending'::user_status
    ELSE 'active'::user_status
END;

ALTER TABLE users
DROP COLUMN pending;

COMMIT;