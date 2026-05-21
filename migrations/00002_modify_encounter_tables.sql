-- +goose Up
ALTER TABLE "encounter" ADD "blizzard_id" bigint NOT NULL;
ALTER TABLE "encounter" ADD "instance_id" bigint;
ALTER TABLE "encounter" ALTER COLUMN "logs_encounter_id" DROP NOT NULL;
ALTER TABLE "encounter" ADD CONSTRAINT "fk_instance_encounters" FOREIGN KEY ("instance_id") REFERENCES "instance"("id");

-- +goose Down
ALTER TABLE "encounter" DROP CONSTRAINT "fk_instance_encounters";
ALTER TABLE "encounter" DROP COLUMN "instance_id";
ALTER TABLE "encounter" DROP COLUMN "blizzard_id";
ALTER TABLE "encounter" ALTER COLUMN "logs_encounter_id" SET NOT NULL;
