CREATE TABLE tokens
(
    id         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    type       VARCHAR(100),
    value      VARCHAR(70) NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    user_id    UUID REFERENCES users (id)
);