-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.subscription
(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    subscriber_id bigint NOT NULL,
    subscribed_to_id bigint NOT NULL,
    created_at date NOT NULL DEFAULT now(),
    updated_at date NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_subscribed_to_fkey FOREIGN KEY (subscribed_to_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT subscriptions_subscriber_fkey FOREIGN KEY (subscriber_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
    CONSTRAINT subscriber_subscribed_to_unique_together UNIQUE(subscriber_id, subscribed_to_id)
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.subscriptions;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
