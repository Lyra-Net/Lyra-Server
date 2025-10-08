CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    display_name TEXT,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    email_encrypted TEXT,
    email_hash TEXT,
    is_2fa BOOLEAN DEFAULT FALSE,
    change_pass_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

ALTER TABLE users
ADD CONSTRAINT unique_username UNIQUE (username),
ADD CONSTRAINT unique_email_hash UNIQUE (email_hash);
CREATE INDEX idx_users_email_hash ON users (email_hash);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    access_jti UUID NOT NULL,
    device_id TEXT,
    browser TEXT,
    os TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE deleted_emails (
    email_hash TEXT PRIMARY KEY,
    deleted_by UUID NOT NULL REFERENCES users(user_id),
    deleted_at TIMESTAMP DEFAULT now(),
    cooldown_until TIMESTAMP, -- for owner deletedby user
    safe_window_until TIMESTAMP, -- for other users
    spam_count INTEGER DEFAULT 0
);

CREATE TABLE trusted_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_id TEXT NOT NULL,
    browser TEXT NOT NULL,
    os TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_trusted_devices_user_device ON trusted_devices (user_id, device_id);
CREATE INDEX idx_trusted_devices_expire ON trusted_devices (expires_at);