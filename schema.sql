BEGIN;
CREATE TABLE requests (
    id serial PRIMARY KEY,
    created_at timestamp with time zone DEFAULT now(),
    method text NOT NULL,
    url text NOT NULL,
    header text NOT NULL,
    body text NOT NULL
);
COMMIT;