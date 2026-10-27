ALTER TABLE teachers
    ADD COLUMN role_id UUID REFERENCES roles(id) ON DELETE RESTRICT;

UPDATE teachers AS t
SET role_id = u.role_id
FROM users AS u
WHERE t.user_id = u.id;

UPDATE teachers
SET role_id = (
    SELECT id
    FROM roles
    WHERE name = 'teacher'
    ORDER BY id
    LIMIT 1
)
WHERE role_id IS NULL;

ALTER TABLE teachers
    ALTER COLUMN role_id SET NOT NULL;

DROP INDEX IF EXISTS users_role_id_idx;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_role_id_fkey,
    DROP COLUMN IF EXISTS role_id;
