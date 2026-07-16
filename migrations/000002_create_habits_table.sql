-- +goose Up
CREATE TABLE IF NOT EXISTS habits (
id TEXT PRIMARY KEY NOT NULL DEFAULT (lower(hex(randomblob(16)))),
user_id TEXT NOT NULL,
name VARCHAR(255) NOT NULL,
description TEXT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);  

CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);

-- +goose Down
DROP TABLE IF EXISTS habits;
DROP INDEX IF EXISTS idx_habits_user_id;