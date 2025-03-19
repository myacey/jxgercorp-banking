BEGIN;

ALTER TABLE users
ADD COLUMN pending boolean DEFAULT false;

UPDATE users
SET pending = CASE
    IF status = 'pending'::user_status THEN true 
    ELSE false
END;

ALTER TABLE users DROP COLUMN status;

DROP TYPE IF EXISTS user_status;

COMMIT;