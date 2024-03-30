-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS personal_dialog (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user1_id BIGINT NOT NULL,
    user2_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user1_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT user1_user2_unique_together UNIQUE(user1_id, user2_id)
);

CREATE OR REPLACE TRIGGER set_timestamp BEFORE
UPDATE ON public.personal_dialog FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

DELETE from public.user;

ALTER TABLE public.personal_message DROP COLUMN IF EXISTS receiver_id;

ALTER TABLE public.personal_message
ADD COLUMN IF NOT EXISTS dialog_id BIGINT;

ALTER TABLE public.personal_message
ADD CONSTRAINT IF NOT EXISTS dialog_id_fk FOREIGN KEY (dialog_id) REFERENCES public.personal_dialog (id) ON UPDATE CASCADE ON DELETE CASCADE;
-- change table users

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
