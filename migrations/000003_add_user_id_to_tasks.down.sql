-- Remove foreign key constraint
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_user;

-- Remove user_id column
ALTER TABLE tasks DROP COLUMN IF EXISTS user_id;
