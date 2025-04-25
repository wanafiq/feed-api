DO
$$
    DECLARE
        user_role_id      UUID;
        moderator_role_id UUID;
        admin_role_id     UUID;
    BEGIN
        -- Fetch role IDs
        SELECT id INTO user_role_id FROM roles WHERE name = 'user';
        SELECT id INTO moderator_role_id FROM roles WHERE name = 'moderator';
        SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';

        -- Insert sample users
        INSERT INTO users (id, username, email, password, is_active, created_at, created_by, updated_at,
                           updated_by, role_id)
        VALUES (gen_random_uuid(), 'user', 'user@gmail.com',
                '$2a$10$b2WFFR5hWO5DfFXNd4xdOuKo.ym5zIJ9T/UDuBEWfYz5AbXFP5TcO', true, CURRENT_TIMESTAMP, 'system',
                CURRENT_TIMESTAMP, 'system', user_role_id),
               (gen_random_uuid(), 'moderator', 'moderator@gmail.com',
                '$2a$10$b2WFFR5hWO5DfFXNd4xdOuKo.ym5zIJ9T/UDuBEWfYz5AbXFP5TcO', true, CURRENT_TIMESTAMP, 'system',
                CURRENT_TIMESTAMP, 'system', moderator_role_id),
               (gen_random_uuid(), 'admin', 'admin@gmail.com',
                '$2a$10$b2WFFR5hWO5DfFXNd4xdOuKo.ym5zIJ9T/UDuBEWfYz5AbXFP5TcO', true, CURRENT_TIMESTAMP, 'system',
                CURRENT_TIMESTAMP, 'system', admin_role_id);
    END
$$;
