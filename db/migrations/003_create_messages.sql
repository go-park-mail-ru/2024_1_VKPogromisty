-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.personal_message
(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id bigint NOT NULL,
    receiver_id bigint NOT NULL,
    content text  NOT NULL DEFAULT ''::text,
    created_at date NOT NULL DEFAULT now(),
    updated_at date NOT NULL DEFAULT now(),
    CONSTRAINT messages_receiver_fkey FOREIGN KEY (receiver_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT messages_sender_fkey FOREIGN KEY (sender_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.messages;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
