-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.sticker (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (author_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.sticker 
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE public.personal_message ADD COLUMN IF NOT EXISTS sticker_id BIGINT DEFAULT NULL;
ALTER TABLE public.personal_message ADD CONSTRAINT fk_sticker_id FOREIGN KEY (sticker_id) REFERENCES public.sticker (id) ON UPDATE CASCADE ON DELETE SET NULL;
---- create above / drop below ----
ALTER TABLE public.personal_message DROP COLUMN IF EXISTS sticker_id;
DROP TRIGGER IF EXISTS set_timestamp ON public.sticker;
DROP TABLE IF EXISTS public.sticker;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
