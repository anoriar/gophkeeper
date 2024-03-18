-- +goose Up
CREATE TABLE IF NOT EXISTS entries (
     id VARCHAR(36) PRIMARY KEY,
     type varchar (20),
     user_id VARCHAR(36) NOT NULL,
     updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
     data bytea NOT NULL,
     meta JSONB NULL
);

CREATE INDEX idx_entry_type ON entries (type);
CREATE INDEX idx_entry_user_id ON entries (user_id);

-- +goose Down
DROP TABLE entries;