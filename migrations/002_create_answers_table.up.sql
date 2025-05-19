CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    user_tg_id BIGINT NOT NULL REFERENCES users (tg_id) ON DELETE CASCADE,
    question_key TEXT NOT NULL,
    short TEXT NOT NULL,
    step INT NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_answer ON answers (
    user_tg_id,
    question_key,
    step
)