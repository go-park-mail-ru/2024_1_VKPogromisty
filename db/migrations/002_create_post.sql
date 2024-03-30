-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.post (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id BIGINT NOT NULL,
    content TEXT NOT NULL DEFAULT ''::TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (author_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE RESTRICT
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.post;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.