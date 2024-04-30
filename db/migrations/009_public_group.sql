-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public_group (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (char_length(name) < 800),
    description TEXT NOT NULL,
    avatar TEXT NOT NULL DEFAULT 'default_avatar.png'::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.public_group
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX public_group_name_trgm_idx ON public.public_group USING gin(name gin_trgm_ops);
---- create above / drop below ----
DROP INDEX IF EXISTS public.public_group_name_trgm_idx;
DROP TRIGGER IF EXISTS set_timestamp ON public.public_group;
DROP TABLE IF EXISTS public_group;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
