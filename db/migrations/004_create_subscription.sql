-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.subscription (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    subscriber_id BIGINT NOT NULL,
    subscribed_to_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT subscriptions_subscribed_to_fkey FOREIGN KEY (subscribed_to_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT subscriptions_subscriber_fkey FOREIGN KEY (subscriber_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT subscriber_subscribed_to_unique_together UNIQUE(subscriber_id, subscribed_to_id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.subscription;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
