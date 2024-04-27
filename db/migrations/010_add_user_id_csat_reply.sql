-- Write your migrate up statements here
ALTER TABLE csat_reply ADD COLUMN user_id BIGINT;
ALTER TABLE csat_reply ADD CONSTRAINT csat_reply_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user (id);
---- create above / drop below ----
ALTER TABLE csat_reply DROP CONSTRAINT csat_reply_user_id_fkey;
ALTER TABLE csat_reply DROP COLUMN user_id;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
