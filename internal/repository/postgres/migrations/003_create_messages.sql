-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.messages
(
    id bigserial NOT NULL,
    sender bigint NOT NULL,
    receiver bigint NOT NULL,
    text text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    creation_date date NOT NULL DEFAULT now(),
    CONSTRAINT messages_pkey PRIMARY KEY (id),
    CONSTRAINT messages_receiver_fkey FOREIGN KEY (receiver)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT messages_sender_fkey FOREIGN KEY (sender)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.messages;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
