
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    max_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



ALTER TABLE students
    ADD COLUMN user_id UUID UNIQUE REFERENCES users(id)
        ON DELETE SET NULL;

ALTER TABLE teachers
    ADD COLUMN user_id UUID UNIQUE REFERENCES users(id)
        ON DELETE SET NULL;


ALTER TABLE applications
    ADD COLUMN user_id UUID REFERENCES users(id)
        ON DELETE SET NULL;

ALTER TABLE applications
    ADD COLUMN organization_id UUID REFERENCES organizations(id)
        ON DELETE RESTRICT;

ALTER TABLE applications
    ADD COLUMN group_id UUID REFERENCES groups(id)
        ON DELETE RESTRICT;


ALTER TABLE applications
    ADD CONSTRAINT applications_type_fields_check
    CHECK (
        (
            type = 'student'
            AND group_id IS NOT NULL
            AND organization_id IS NOT NULL
        )
        OR
        (
            type = 'teacher'
            AND group_id IS NULL
            AND organization_id IS NOT NULL
        )
    );


ALTER TABLE groups
    ADD CONSTRAINT groups_id_organization_unique
    UNIQUE (id, organization_id);

ALTER TABLE applications
    ADD CONSTRAINT applications_group_organization_fk
    FOREIGN KEY (group_id, organization_id)
    REFERENCES groups(id, organization_id)
    ON DELETE RESTRICT;

CREATE INDEX applications_organization_status_idx
    ON applications (organization_id, status);

CREATE INDEX applications_user_id_idx
    ON applications (user_id);

CREATE INDEX students_group_id_idx
    ON students (group_id);

CREATE INDEX teachers_organization_id_idx
    ON teachers (organization_id);