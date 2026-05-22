-- +goose Up
ALTER TABLE "instance" ALTER COLUMN "name" TYPE jsonb USING "name"::jsonb;
ALTER TABLE "encounter" ALTER COLUMN "name" TYPE jsonb USING "name"::jsonb;
ALTER TABLE "encounter" DROP COLUMN "names";

-- +goose Down
ALTER TABLE "instance" ALTER COLUMN "name" TYPE text USING "name"::jsonb;
ALTER TABLE "encounter" ALTER COLUMN "name" TYPE text USING "name"::jsonb;
ALTER TABLE "encounter" ADD COLUMN "names" jsonb;
