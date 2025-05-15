CREATE TABLE IF NOT EXISTS users (
    tg_id BIGINT PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    username TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    remind_stage INT DEFAULT 0 CHECK (remind_stage IN (0, 1, 2)),
    remind_at TIMESTAMP NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    current_step INT NOT NULL DEFAULT 0,
    max_step_reached INT NOT NULL DEFAULT 0
);