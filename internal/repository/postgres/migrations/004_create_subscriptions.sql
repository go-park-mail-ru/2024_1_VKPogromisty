-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.subscriptions
(
    id bigserial NOT NULL,
    subscriber bigint NOT NULL,
    subscribed_to bigint NOT NULL,
    creation_date date NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_pkey PRIMARY KEY (id),
    CONSTRAINT subscriptions_subscribed_to_fkey FOREIGN KEY (subscribed_to)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT subscriptions_subscriber_fkey FOREIGN KEY (subscriber)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.subscriptions;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
