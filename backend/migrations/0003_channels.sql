-- +goose Up
CREATE TABLE IF NOT EXISTS channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL CHECK (type IN ('private', 'group', 'channel')),
    name TEXT,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    members UUID[] NOT NULL DEFAULT '{}',
    permissions JSONB DEFAULT '{}'::jsonb,
    security_label TEXT NOT NULL DEFAULT 'public' CHECK (security_label IN ('public', 'internal', 'confidential')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for finding user's channels
CREATE INDEX IF NOT EXISTS idx_channels_members ON channels USING GIN(members);

-- Index for owner lookup
CREATE INDEX IF NOT EXISTS idx_channels_owner ON channels(owner_id);

-- +goose Down
DROP INDEX IF EXISTS idx_channels_members;
DROP INDEX IF EXISTS idx_channels_owner;
DROP TABLE IF EXISTS channels;
