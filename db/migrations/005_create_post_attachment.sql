-- Write your migrate up statements here
select * FROM public.user;
-- CREATE TABLE IF NOT EXISTS public.post_attachment (
--     id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     post_id BIGINT NOT NULL,
--     file_name TEXT UNIQUE NOT NULL,
--     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--     updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--     FOREIGN KEY (post_id) REFERENCES public.post (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
---- create above / drop below ----
DROP TABLE IF EXISTS public.post_attachment;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
