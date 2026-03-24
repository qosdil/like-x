CREATE TABLE public.users (
    id serial NOT NULL,

    -- Designed for Nano ID (random, no timestamp, shorter than UUID)
    public_id text unique NOT NULL,

    password_hash VARCHAR(255),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
