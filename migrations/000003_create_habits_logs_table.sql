-- +goose Up
CREATE TABLE IF NOT EXISTS habits_logs (
id TEXT PRIMARY KEY NOT NULL DEFAULT (lower(hex(randomblob(16)))),
habit_id TEXT NOT NULL,
completed_date DATE NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
UNIQUE (habit_id, completed_date)
);

CREATE INDEX IF NOT EXISTS idx_habits_logs_habit_id ON habits_logs(habit_id);
CreATE INDEX IF NOT EXISTS idx_habits_logs_completed_date ON habits_logs(completed_date);

-- +goose Down
DROP TABLE IF EXISTS habits_logs;
DROP INDEX IF EXISTS idx_habits_logs_habit_id;
DROP INDEX IF EXISTS idx_habits_logs_completed_date;