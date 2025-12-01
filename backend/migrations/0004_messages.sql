-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    content BYTEA NOT NULL, -- Encrypted content
    content_type TEXT NOT NULL CHECK (content_type IN ('text', 'image', 'audio', 'video')),
    encryption_meta JSONB NOT NULL DEFAULT '{}'::jsonb, -- IV, algorithm info, key hints
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for channel message lookups (most common query)
CREATE INDEX IF NOT EXISTS idx_messages_channel_timestamp ON messages(channel_id, timestamp DESC) WHERE deleted = false;

-- Index for sender lookups
CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender_id);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_channel_timestamp;
DROP INDEX IF EXISTS idx_messages_sender;
DROP TABLE IF EXISTS messages;
