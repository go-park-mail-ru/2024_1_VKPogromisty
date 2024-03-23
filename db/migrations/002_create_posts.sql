-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.post
(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id bigint NOT NULL,
    content text  NOT NULL DEFAULT ''::text,
    created_at date NOT NULL DEFAULT now(),
    updated_at date NOT NULL DEFAULT now(),
    CONSTRAINT FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.posts;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
