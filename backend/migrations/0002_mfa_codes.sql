CREATE TABLE IF NOT EXISTS mfa_codes (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    code        TEXT NOT NULL,
    expires_at  TIMESTAMP NOT NULL
);
