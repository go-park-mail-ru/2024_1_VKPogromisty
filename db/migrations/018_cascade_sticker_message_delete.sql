-- Write your migrate up statements here
ALTER TABLE public.personal_message DROP CONSTRAINT fk_sticker_id;

ALTER TABLE public.personal_message ADD CONSTRAINT fk_sticker_id FOREIGN KEY (sticker_id) REFERENCES public.sticker (id) ON UPDATE CASCADE ON DELETE CASCADE;
---- create above / drop below ----
ALTER TABLE public.personal_message DROP CONSTRAINT fk_sticker_id;

ALTER TABLE public.personal_message ADD CONSTRAINT fk_sticker_id FOREIGN KEY (sticker_id) REFERENCES public.sticker (id) ON UPDATE CASCADE ON DELETE SET NULL;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
