-- Write your migrate up statements here
ALTER TABLE subscriptions
ADD CONSTRAINT subscriber_subscribed_to_unique_together
UNIQUE(subscriber, subscribed_to);
---- create above / drop below ----
ALTER TABLE subscriptions
DROP CONSTRAINT subscriber_subscribed_to_unique_together;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
