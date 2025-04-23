CREATE TABLE roles
(
    id          UUID    DEFAULT gen_random_uuid() PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    level       INT          NOT NULL,
    description TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL,
    created_by  VARCHAR(100) NOT NULL,
    updated_at  TIMESTAMPTZ,
    updated_by  VARCHAR(100)
);