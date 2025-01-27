-- Create schema_migrations table if it doesn't exist
CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL,
    CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);
