-- +goose Up
CREATE TYPE payment_method AS ENUM (
    'UNSPECIFIED',
    'CARD',
    'SBP',
    'CREDIT_CARD',
    'INVESTOR_MONEY'
);

CREATE TYPE order_status AS ENUM (
    'PENDING_PAYMENT',
    'PAID',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS orders (
                        id SERIAL PRIMARY KEY,
                        order_uuid VARCHAR(36) NOT NULL UNIQUE,
                        user_uuid VARCHAR(36) NOT NULL,
                        part_uuids VARCHAR(36)[] NOT NULL,
                        total_price DECIMAL(10, 2) NOT NULL,
                        transaction_uuid VARCHAR(36),
                        payment_method payment_method,
                        status order_status NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_orders_user_uuid ON orders(user_uuid);
CREATE INDEX idx_orders_status ON orders(status);

-- +goose Down
DROP TABLE orders;
DROP TYPE payment_method;
DROP TYPE order_status;
