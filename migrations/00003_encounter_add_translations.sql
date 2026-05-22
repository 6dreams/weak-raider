-- +goose Up
ALTER TABLE "encounter" ADD "names" jsonb NOT NULL;

-- +goose Down
ALTER TABLE "encounter" DROP "names";
