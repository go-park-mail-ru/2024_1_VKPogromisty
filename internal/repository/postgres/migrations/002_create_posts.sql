-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.posts
(
    id bigserial NOT NULL,
    author_id bigint NOT NULL,
    text text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    creation_date date NOT NULL DEFAULT now(),
    CONSTRAINT posts_pkey PRIMARY KEY (id),
    CONSTRAINT posts_author_id_fkey FOREIGN KEY (author_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.posts;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
