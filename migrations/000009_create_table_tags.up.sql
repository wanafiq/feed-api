CREATE TABLE tags
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL
)