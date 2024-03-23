-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.post_attachment
(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    file_name text UNIQUE NOT NULL,
    post_id bigint NOT NULL,
    created_at date NOT NULL DEFAULT now(),
    updated_at date NOT NULL DEFAULT now(),
    FOREIGN KEY (post_id)
        REFERENCES public.posts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.post_attachments;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
