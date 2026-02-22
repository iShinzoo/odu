CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    amount NUMERIC NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_status_logs (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);