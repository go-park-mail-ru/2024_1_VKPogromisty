-- Write your migrate up statements here
ALTER TABLE csat_pool DROP CONSTRAINT csat_pool_author_id_fkey;
ALTER TABLE csat_pool DROP COLUMN author_id;
---- create above / drop below ----
ALTER TABLE csat_pool ADD COLUMN author_id BIGINT;
ALTER TABLE csat_pool ADD CONSTRAINT csat_pool_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.admin (id) ON UPDATE CASCADE ON DELETE RESTRICT;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
