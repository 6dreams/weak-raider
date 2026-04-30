-- +goose Up
CREATE TABLE IF NOT EXISTS guild(
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  wowaudit_key TEXT NOT NULL,
  updated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS champion (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  class text NOT NULL,
  realm text NOT NULL,
  note text DEFAULT '',
  guild_id BIGINT,
  blizzard_id BIGINT NOT NULL,
  updated TIMESTAMP,
  FOREIGN KEY (guild_id) REFERENCES guild(id)
);
-- +goose Down
SELECT 'DROP TABLE IF EXISTS guild';
