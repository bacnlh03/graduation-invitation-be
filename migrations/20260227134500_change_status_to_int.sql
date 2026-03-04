-- +goose Up
-- +goose StatementBegin
ALTER TABLE guests ALTER COLUMN status DROP DEFAULT;
ALTER TABLE guests ALTER COLUMN status TYPE SMALLINT USING (CASE WHEN status THEN 2 ELSE 0 END);
ALTER TABLE guests ALTER COLUMN status SET DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE guests ALTER COLUMN status DROP DEFAULT;
ALTER TABLE guests ALTER COLUMN status TYPE BOOLEAN USING (CASE WHEN status = 2 THEN TRUE ELSE FALSE END);
ALTER TABLE guests ALTER COLUMN status SET DEFAULT FALSE;
-- +goose StatementEnd
