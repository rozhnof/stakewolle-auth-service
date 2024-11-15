CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    referrer_id UUID REFERENCES users(id),
    hash_password TEXT NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE referral_code (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    code TEXT NOT NULL UNIQUE,
    expired_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_refresh_token ON session (refresh_token);
