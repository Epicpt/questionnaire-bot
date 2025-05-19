CREATE TABLE IF NOT EXISTS emails (
    id SERIAL PRIMARY KEY,
    user_tg_id BIGINT NOT NULL UNIQUE REFERENCES users (tg_id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK (
        status IN ('pending', 'sent', 'failed')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ
);