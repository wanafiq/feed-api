CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- Enable UUID generation support

CREATE TABLE users
(
    id         UUID    DEFAULT gen_random_uuid() PRIMARY KEY,
    username   VARCHAR(100) UNIQUE NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(70)        NOT NULL,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ         NOT NULL,
    created_by VARCHAR(100)        NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100)
);