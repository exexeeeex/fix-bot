DROP INDEX IF EXISTS teachers_organization_id_idx;
DROP INDEX IF EXISTS students_group_id_idx;
DROP INDEX IF EXISTS applications_user_id_idx;
DROP INDEX IF EXISTS applications_organization_status_idx;

ALTER TABLE applications
    DROP CONSTRAINT IF EXISTS applications_group_organization_fk;

ALTER TABLE applications
    DROP CONSTRAINT IF EXISTS applications_type_fields_check;

ALTER TABLE groups
    DROP CONSTRAINT IF EXISTS groups_id_organization_unique;

ALTER TABLE applications
    DROP COLUMN IF EXISTS group_id,
    DROP COLUMN IF EXISTS organization_id,
    DROP COLUMN IF EXISTS user_id;

ALTER TABLE teachers
    DROP COLUMN IF EXISTS user_id;

ALTER TABLE students
    DROP COLUMN IF EXISTS user_id;

DROP TABLE IF EXISTS users;