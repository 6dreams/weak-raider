-- +goose Up
CREATE TABLE "guild" ("id" bigserial NOT NULL,"name" text NOT NULL,"wow_audit_key" text NOT NULL,"updated_at" timestamptz NOT NULL DEFAULT now(),PRIMARY KEY ("id"));
CREATE TABLE "character" ("id" bigserial NOT NULL,"name" text NOT NULL,"class" text NOT NULL,"rank" text NOT NULL,"role" text NOT NULL,"realm" text NOT NULL,"note" text NOT NULL DEFAULT '',"guild_id" bigint,"updated_at" timestamptz NOT NULL DEFAULT now(),PRIMARY KEY ("id"),CONSTRAINT "fk_character_guild" FOREIGN KEY ("guild_id") REFERENCES "guild"("id"));
CREATE TABLE "encounter" ("id" bigserial NOT NULL,"name" text NOT NULL,"logs_encounter_id" bigint NOT NULL,"updated_at" timestamptz NOT NULL DEFAULT now(),PRIMARY KEY ("id"));
CREATE TABLE "season" ("id" bigserial NOT NULL,"blizzard_id" bigint NOT NULL,"wow_audit_id" bigint,"name" text NOT NULL,"is_current" boolean NOT NULL,"updated_at" timestamptz NOT NULL DEFAULT now(),PRIMARY KEY ("id"));
CREATE TABLE "instance" ("id" bigserial NOT NULL,"season_id" bigint,"name" text NOT NULL,"is_raid" boolean NOT NULL,"updated_at" timestamptz NOT NULL DEFAULT now(),PRIMARY KEY ("id"),CONSTRAINT "fk_instance_season" FOREIGN KEY ("season_id") REFERENCES "season"("id"));
CREATE TABLE "slot" ("id" bigserial NOT NULL,"name" text NOT NULL,PRIMARY KEY ("id"));
CREATE TABLE "item" ("id" bigserial NOT NULL,"name" text NOT NULL,"is_unique" boolean NOT NULL,"sockets" bigint,"instance_id" bigint,"encounter_id" bigint,"slot_id" bigint,PRIMARY KEY ("id"),CONSTRAINT "fk_item_instance" FOREIGN KEY ("instance_id") REFERENCES "instance"("id"),CONSTRAINT "fk_item_encounter" FOREIGN KEY ("encounter_id") REFERENCES "encounter"("id"),CONSTRAINT "fk_item_slot" FOREIGN KEY ("slot_id") REFERENCES "slot"("id"));

-- +goose Down
DROP TABLE IF EXISTS "item";
DROP TABLE IF EXISTS "slot";
DROP TABLE IF EXISTS "instance";
DROP TABLE IF EXISTS "season";
DROP TABLE IF EXISTS "encounter";
DROP TABLE IF EXISTS "character";
DROP TABLE IF EXISTS guild;
