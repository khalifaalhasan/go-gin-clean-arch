CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    CHECK (stock >= 0)
); -- <--- JANGAN LUPA TITIK KOMA INI (Pemisah Tabel)

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE idempotency_keys(
    key VARCHAR NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    response_code INT,
    response_body JSONB, -- buat nyimpen respon sukses 
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
)

