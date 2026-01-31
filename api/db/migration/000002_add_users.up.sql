CREATE TABLE users (
  id bigserial PRIMARY KEY,
  username varchar NOT NULL UNIQUE,
  password varchar NOT NULL,
  full_name varchar NOT NULL,
  email varchar NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT (now())
);

-- (Opsional) Bikin Index biar pencarian cepat
CREATE INDEX ON users (username);
CREATE INDEX ON users (email);