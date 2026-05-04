-- +goose Up
CREATE TABLE "guild" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "wow_audit_key" 
  text NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"));

CREATE TABLE "character" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "class" text NOT NULL,
  "rank" text NOT NULL,
  "role" text NOT NULL,
  "realm" text NOT NULL,
  "note" text NOT NULL DEFAULT '',
  "guild_id" bigserial,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_character_guild" FOREIGN KEY ("guild_id") REFERENCES "guild"("id"));
-- +goose Down
DROP TABLE IF EXISTS character;
DROP TABLE IF EXISTS guild;

