-- migrate:up
CREATE TABLE IF NOT EXISTS tokens (
  "userId" BIGINT PRIMARY KEY,
  "token" TEXT NOT NULL
);

-- migrate:down
