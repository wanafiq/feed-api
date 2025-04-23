INSERT INTO roles (id, name, level, description, is_active, created_at, created_by, updated_at, updated_by)
VALUES (gen_random_uuid(), 'user', 1, 'Regular user', true, CURRENT_TIMESTAMP, 'system', CURRENT_TIMESTAMP, 'system'),
       (gen_random_uuid(), 'moderator', 2, 'Moderator with extended privileges', true, CURRENT_TIMESTAMP, 'system',
        CURRENT_TIMESTAMP, 'system'),
       (gen_random_uuid(), 'admin', 3, 'Administrator with full access', true, CURRENT_TIMESTAMP, 'system',
        CURRENT_TIMESTAMP, 'system');
