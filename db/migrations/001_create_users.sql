-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.user
(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    salt TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL CHECK (),
    created_at date NOT NULL DEFAULT now(),
    updated_at date NOT NULL DEFAULT now(),
    avatar text  NOT NULL DEFAULT 'default_avatar.png'::text,
    date_of_birth date NOT NULL,
);
---- create above / drop below ----
DROP TABLE IF EXISTS public.users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
