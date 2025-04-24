CREATE TABLE posts
(
    id           UUID    DEFAULT gen_random_uuid() PRIMARY KEY,
    author_id    UUID REFERENCES users (id),
    title        TEXT         NOT NULL,
    slug         TEXT         NOT NULL,
    content      TEXT         NOT NULL,
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ  NOT NULL,
    created_by   VARCHAR(100) NOT NULL,
    updated_at   TIMESTAMPTZ,
    updated_by   VARCHAR(100)
)