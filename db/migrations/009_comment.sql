-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.comment (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    content TEXT NOT NULL DEFAULT ''::TEXT,
    attachments TEXT[] DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (author_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES public.post (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.comment 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();
---- create above / drop below ----
DROP TABLE IF EXISTS public.comment;
DROP TRIGGER IF EXISTS set_timestamp ON public.comment;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
