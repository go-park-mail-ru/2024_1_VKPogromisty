-- Write your migrate up statements here
CREATE TABLE admin (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES public.user (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE csat_pool (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT,
    author_id BIGINT,
    is_active BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (author_id) REFERENCES public.admin (id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE csat_question (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pool_id BIGINT,
    question TEXT,
    worst_case TEXT,
    best_case TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (pool_id) REFERENCES public.csat_pool (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE csat_reply (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    question_id BIGINT,
    score INT CHECK (score >= 1 AND score <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (question_id) REFERENCES public.csat_question (id) ON UPDATE CASCADE ON DELETE CASCADE
);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
DROP TABLE IF EXISTS csat_reply;
DROP TABLE IF EXISTS csat_question;
DROP TABLE IF EXISTS csat_pool;
DROP TABLE IF EXISTS admin;
