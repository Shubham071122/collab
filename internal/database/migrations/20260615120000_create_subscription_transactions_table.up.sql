CREATE TABLE subscription_transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id VARCHAR(255) NOT NULL,
    payment_id VARCHAR(255),
    amount INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    billing_reason VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
