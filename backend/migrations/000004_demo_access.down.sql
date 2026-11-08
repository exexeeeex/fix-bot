DELETE FROM visits
WHERE student_id IN (
    SELECT students.id
    FROM students
    JOIN groups ON groups.id = students.group_id
    WHERE groups.organization_id = '20000000-0000-0000-0000-000000000002'
);

DELETE FROM applications
WHERE organization_id = '20000000-0000-0000-0000-000000000002';

DELETE FROM students
WHERE group_id IN (
    SELECT id
    FROM groups
    WHERE organization_id = '20000000-0000-0000-0000-000000000002'
);

DELETE FROM teachers
WHERE organization_id = '20000000-0000-0000-0000-000000000002';

DELETE FROM users
WHERE id IN (
    '50000000-0000-0000-0000-000000000001',
    '50000000-0000-0000-0000-000000000002'
);

DELETE FROM groups
WHERE organization_id = '20000000-0000-0000-0000-000000000002';

DELETE FROM organizations
WHERE id = '20000000-0000-0000-0000-000000000002';
