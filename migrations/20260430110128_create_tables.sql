-- +goose Up
CREATE TABLE IF NOT EXISTS guild(
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  wowaudit_key TEXT NOT NULL,
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "character" (
  id BIGSERIAL PRIMARY KEY,
  class text NOT NULL,
  name text NOT NULL,
  note text DEFAULT '',
  rank text NOT NULL,
  realm text NOT NULL,
  role text NOT NULL,
  guild_id BIGINT,
  updated_at TIMESTAMP,
  wowaudit_updated_at TIMESTAMP,

  FOREIGN KEY (guild_id) REFERENCES guild(id)
);
-- +goose Down
DROP TABLE IF EXISTS character;
DROP TABLE IF EXISTS guild;

