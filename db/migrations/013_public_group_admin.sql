-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public_group_admin (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES public.user(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES public.public_group(id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE (user_id, group_id)
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.public_group_admin
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
---- create above / drop below ----
DROP TRIGGER IF EXISTS set_timestamp ON public.public_group_admin;
DROP TABLE IF EXISTS public.public_group_admin;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
