-- migrate:up
CREATE TABLE IF NOT EXISTS tokens (
  "id" BIGINT PRIMARY KEY,
  "token" TEXT NOT NULL
);

-- migrate:down
