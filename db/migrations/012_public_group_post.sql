-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public_group_post (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_group_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (public_group_id) REFERENCES public.public_group(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES public.post (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT public_group_post_post_unique UNIQUE (post_id)
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.public_group_post
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
---- create above / drop below ----
DROP TRIGGER IF EXISTS set_timestamp ON public.public_group_post;
DROP TABLE IF EXISTS public.public_group_post;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
