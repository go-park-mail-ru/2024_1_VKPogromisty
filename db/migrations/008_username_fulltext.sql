-- Write your migrate up statements here
CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE public.user ADD COLUMN full_name TEXT;

UPDATE public.user SET full_name = first_name || ' ' || last_name;

CREATE INDEX user_full_name_trgm_idx ON public.user USING gin(full_name gin_trgm_ops);
---- create above / drop below ----
DROP INDEX IF EXISTS public.user_full_name_trgm_idx;

ALTER TABLE public.user DROP COLUMN IF EXISTS full_name;

DROP EXTENSION IF EXISTS pg_trgm;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
