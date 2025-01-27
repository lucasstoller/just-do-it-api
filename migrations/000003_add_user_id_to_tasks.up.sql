-- Add user_id column
ALTER TABLE tasks ADD COLUMN user_id INTEGER NOT NULL DEFAULT 1;

-- Add foreign key constraint
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE;
