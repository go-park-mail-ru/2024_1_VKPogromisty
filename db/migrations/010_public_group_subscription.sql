-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public_group_subscription (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_group_id BIGINT NOT NULL,
    subscriber_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (public_group_id) REFERENCES public.public_group(id) ON DELETE CASCADE,
    FOREIGN KEY (subscriber_id) REFERENCES public.user(id) ON DELETE CASCADE,
    CONSTRAINT subscriber_public_group_unique_together UNIQUE(subscriber_id, public_group_id)
);

CREATE OR REPLACE TRIGGER set_timestamp
BEFORE UPDATE ON public.public_group_subscription
FOR EACH ROW 
EXECUTE PROCEDURE trigger_set_timestamp();
---- create above / drop below ----
DROP TRIGGER IF EXISTS set_timestamp ON public.public_group_subscription;
DROP TABLE IF EXISTS public_group_subscription;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
