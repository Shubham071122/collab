CREATE TYPE subscription_tier AS ENUM ('free', 'silver', 'gold');
CREATE TYPE subscription_status AS ENUM ('active', 'past_due', 'canceled', 'incomplete');
CREATE TYPE payment_provider AS ENUM ('none', 'razorpay', 'stripe');

CREATE TABLE user_subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    tier subscription_tier NOT NULL DEFAULT 'free',
    status subscription_status NOT NULL DEFAULT 'active',
    provider payment_provider NOT NULL DEFAULT 'none',
    provider_subscription_id VARCHAR(255),
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
