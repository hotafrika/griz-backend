-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS codes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    link VARCHAR NOT NULL,
    hash VARCHAR NULL UNIQUE ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE);

CREATE INDEX IF NOT EXISTS idx_codes_hash ON codes(hash);

CREATE TRIGGER IF NOT EXISTS on_update_code_update_time AFTER UPDATE ON codes FOR EACH ROW BEGIN
    UPDATE codes SET updated_at = CURRENT_TIMESTAMP WHERE id = old.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE codes;
-- +goose StatementEnd
