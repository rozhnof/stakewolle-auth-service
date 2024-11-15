CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id REFERENCES users(id)
    username VARCHAR(50) NOT NULL UNIQUE,
    hash_password TEXT NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE referral_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    expired_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    refresh_token VARCHAR(255),
    expired_at TIMESTAMP,
    is_revoked BOOLEAN,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_refresh_token ON session (refresh_token);
