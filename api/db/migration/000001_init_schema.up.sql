CREATE TABLE comments (
    id bigserial PRIMARY KEY,
    username varchar(30) NOT NULL,
    content text NOT NULL,
    ip_address varchar NOT NULL,
    is_hidden boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE INDEX ON comments (username);
CREATE INDEX ON comments (created_at);