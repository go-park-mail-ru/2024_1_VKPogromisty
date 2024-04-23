-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS post_like (
  id SERIAL PRIMARY KEY,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  FOREIGN KEY (post_id) REFERENCES public.post (id) ON UPDATE CASCADE ON DELETE CASCADE
  FOREIGN KEY (user_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE
  CONSTRAINT post_like_post_user_unique UNIQUE (post_id, user_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS public.post_like;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
