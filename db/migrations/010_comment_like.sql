-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS comment_like (
  id SERIAL PRIMARY KEY,
  comment_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  FOREIGN KEY (comment_id) REFERENCES public.comment (id) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT comment_like_comment_user_unique UNIQUE (comment_id, user_id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.comment_like;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
