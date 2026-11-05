-- +goose Up
TRUNCATE TABLE "item", "encounter", "instance", "season";

ALTER TABLE "season" ADD COLUMN "expansion_id" bigint NOT NULL;

-- +goose Down
ALTER TABLE "season" DROP COLUMN "expansion_id";
