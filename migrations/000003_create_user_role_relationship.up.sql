ALTER TABLE users
    ADD COLUMN role_id UUID,
    ADD CONSTRAINT fk_user_role FOREIGN KEY (role_id)
        REFERENCES roles (id) ON DELETE SET NULL; -- if role is deleted, user's role_id will be set to null