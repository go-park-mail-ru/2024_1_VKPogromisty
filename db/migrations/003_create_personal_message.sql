-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.personal_message (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,
    content TEXT NOT NULL DEFAULT ''::TEXT,
    attachments TEXT[] DEFAULT ARRAY[]::TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT messages_receiver_fkey FOREIGN KEY (receiver_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT messages_sender_fkey FOREIGN KEY (sender_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE RESTRICT
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.message;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.