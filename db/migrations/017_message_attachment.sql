-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.message_attachment (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message_id BIGINT NOT NULL,
    file_name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (message_id) REFERENCES public.personal_message (id) ON UPDATE CASCADE ON DELETE CASCADE
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.message_attachment;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
