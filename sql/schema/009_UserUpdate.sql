-- +goose Up
CREATE TYPE user_role AS ENUM ('member', 'moderator', 'admin');

-- Convert the existing role column, casting strings to the enum type
ALTER TABLE users ALTER COLUMN role TYPE user_role USING role::user_role;
-- Apply defaults and constraints
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'member';
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

-- +goose Down
ALTER TABLE users ALTER COLUMN role TYPE TEXT;
ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role DROP NOT NULL;
DROP TYPE user_role CASCADE;
