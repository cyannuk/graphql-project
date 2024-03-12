-- migrate:up
ALTER TABLE users ADD COLUMN "role" INT NOT NULL DEFAULT 0;

-- migrate:down
