-- +goose Up
ALTER TABLE entries ADD COLUMN original_id VARCHAR(255) NOT NULL DEFAULT '';
CREATE UNIQUE INDEX uniq_ids_entries_original_id_user_id ON entries (original_id, user_id);




-- +goose Down
DROP INDEX IF EXISTS uniq_ids_entries_original_id_user_id;
ALTER TABLE entries DROP COLUMN original_id;