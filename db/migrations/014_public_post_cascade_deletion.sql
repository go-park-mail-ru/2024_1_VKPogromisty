-- Write your migrate up statements here
CREATE OR REPLACE FUNCTION delete_related_public_group_posts() RETURNS TRIGGER AS $$
BEGIN
  DELETE FROM public.post
  WHERE public.post.id = OLD.post_id;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER delete_related_public_group_posts
AFTER DELETE ON public.public_group_post
FOR EACH ROW EXECUTE FUNCTION delete_related_public_group_posts();
---- create above / drop below ----
DROP FUNCTION IF EXISTS delete_related_public_group_posts() CASCADE;
DROP TRIGGER IF EXISTS delete_related_public_group_posts ON public.public_group_post;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
