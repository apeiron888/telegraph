-- +goose Up
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Nullable for failed login attempts
    action TEXT NOT NULL,
    resource TEXT,
    ip_address TEXT,
    result TEXT NOT NULL CHECK (result IN ('success', 'failure')),
    details TEXT,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for user audit trail lookups
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_timestamp ON audit_logs(user_id, timestamp DESC);

-- Index for action-based queries
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- +goose Down
DROP INDEX IF EXISTS idx_audit_logs_user_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP TABLE IF EXISTS audit_logs;
