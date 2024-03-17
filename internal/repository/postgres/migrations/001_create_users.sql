-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS public.users
(
    id BIGSERIAL NOT NULL,
    first_name text COLLATE pg_catalog."default" NOT NULL,
    last_name text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    salt text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    registration_date date NOT NULL DEFAULT now(),
    avatar text COLLATE pg_catalog."default" NOT NULL DEFAULT 'default_avatar.png'::text,
    date_of_birth date NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email)
)
---- create above / drop below ----
DROP TABLE IF EXISTS public.users;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
