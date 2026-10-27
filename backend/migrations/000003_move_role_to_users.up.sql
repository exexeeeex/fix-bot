DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM teachers
        WHERE user_id IS NULL
    ) THEN
        RAISE EXCEPTION
            'cannot move teacher roles to users: every teacher must be linked to a user';
    END IF;
END;
$$;

ALTER TABLE users
    ADD COLUMN role_id UUID;

UPDATE users AS u
SET role_id = t.role_id
FROM teachers AS t
WHERE t.user_id = u.id;

ALTER TABLE users
    ADD CONSTRAINT users_role_id_fkey
    FOREIGN KEY (role_id) REFERENCES roles(id)
    ON DELETE RESTRICT;

CREATE INDEX users_role_id_idx
    ON users (role_id);

ALTER TABLE teachers
    DROP COLUMN role_id;
