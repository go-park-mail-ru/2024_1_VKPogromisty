-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.user (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- longest personal name in the world is 746 characters
    first_name TEXT NOT NULL CHECK (char_length(first_name) < 800),
    last_name TEXT NOT NULL CHECK (char_length(first_name) < 800),
    hashed_password TEXT NOT NULL,
    salt TEXT NOT NULL,
    -- longest email address is 320 characters
    email TEXT UNIQUE NOT NULL CHECK (char_length(email) < 350),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    avatar TEXT NOT NULL DEFAULT 'default_avatar.png'::TEXT,
    date_of_birth DATE NOT NULL
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.user;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.